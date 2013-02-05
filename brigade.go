package brigade

import (
	"container/list"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

var Errors *list.List
var ScanDirs *list.List

var CopyFiles chan string
var DeleteFiles chan string

type S3Connection struct {
	Source       *s3.S3
	Dest         *s3.S3
	SourceBucket *s3.Bucket
	DestBucket   *s3.Bucket
}

func S3Connect(t *Target) *s3.S3 {
	auth := aws.Auth{t.AccessKey, t.SecretAccessKey}
	return s3.New(auth, aws.Region{S3Endpoint: t.Server})
}

func S3Init() *S3Connection {
	s := &S3Connection{S3Connect(Config.Source), S3Connect(Config.Dest), nil, nil}

	if s.Source == nil {
		log.Fatalf("Could not connect to S3 endpoint %s", Config.Source.Server)
	}

	if s.Dest == nil {
		log.Fatalf("Could not connect to S3 endpoint %s", Config.Dest.Server)
	}

	s.SourceBucket = s.Source.Bucket(Config.Source.BucketName)
	s.DestBucket = s.Source.Bucket(Config.Dest.BucketName)

	return s
}

func (s *S3Connection) fileCopier() {
	// pull files off channel, copy with permissions
  

}

var worker []*S3Connection

func Init() {
	ScanDirs = list.New()
	Errors = list.New()

	CopyFiles = make(chan string, 1000)
	DeleteFiles = make(chan string, 100)

  worker = make([]*S3Connection, Config.Workers)

	// spawn workers
	for i := 0; i < Config.Workers; i++ {
    worker[i] = S3Init()
		go worker[i].fileCopier()
	}
}

func (s *S3Connection) CopyBucket() {
	ScanDirs.PushBack("")
	for ScanDirs.Len() > 0 {
		dir, ok := ScanDirs.Remove(ScanDirs.Front()).(string)
		if !ok {
			log.Fatalf("Invalid value found on directory queue")
		}
		err := s.CopyDirectory(dir)
		if err != nil {
			Errors.PushBack(err)
		}
	}
}

func inList(input string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if input == list[i] {
			return true
		}
	}
	return false
}

var nilKey s3.Key

func findKey(name string, list *s3.ListResp) (s3.Key, bool) {
	for i := 0; i < len(list.Contents); i++ {
		if list.Contents[i].Key == name {
			return list.Contents[i], true
		}
	}
	return nilKey, false
}

func keyChanged(src s3.Key, dest s3.Key) bool {
	return src.Size != dest.Size || src.ETag != dest.ETag || src.StorageClass != dest.StorageClass
}

func (s *S3Connection) CopyDirectory(dir string) error {

	sourceList, err := s.SourceBucket.List(dir, "/", "", 1000)
	if err != nil {
		return err
	}

	destList, err := s.DestBucket.List(dir, "/", "", 1000)
	if err != nil {
		return err
	}

	// push subdirectories onto directory queue
	for i := 0; i < len(sourceList.CommonPrefixes); i++ {
		ScanDirs.PushBack(sourceList.CommonPrefixes[i])
	}

	// push subdirectories that no longer exist onto delete queue
	for i := 0; i < len(destList.CommonPrefixes); i++ {
		if !inList(destList.CommonPrefixes[i], sourceList.CommonPrefixes) {
			ScanDirs.PushBack(destList.CommonPrefixes[i])
		}
	}

	// push changed files onto file queue
	for i := 0; i < len(sourceList.Contents); i++ {
		key := sourceList.Contents[i]
		existing, found := findKey(key.Key, destList)
		if !found || keyChanged(key, existing) {
			CopyFiles <- key.Key
		}
	}

	// push removed files onto delete list
	for i := 0; i < len(destList.Contents); i++ {
		key := destList.Contents[i]
		_, found := findKey(key.Key, sourceList)
		if !found {
			DeleteFiles <- key.Key
		}
	}

	return nil
}


