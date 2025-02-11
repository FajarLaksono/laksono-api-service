// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package main

// TODO: Implement the function, Need to upgrage the GO version and adjust the code to use the Go packages.

// import (
// 	"os"

// 	"fajarlaksono.github.io/laksono-api-service/app/config"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/aws/awserr"
// 	"github.com/aws/aws-sdk-go-v2/aws/session"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/caarlos0/env"
// 	"github.com/pkg/errors"
// 	"github.com/sirupsen/logrus"
// 	log "github.com/sirupsen/logrus"
// )

// func main() {
// 	log.SetOutput(os.Stdout)
// 	log.Info("Starting S3 Migration Tool")

// 	// Load configuration
// 	conf := &config.Config{}

// 	if err := env.Parse(conf); err != nil {
// 		log.WithError(err).Fatal("Failed to load environment variables")

// 		return
// 	}

// 	log.Infof("region %s, service endpoint : %s, bucket : %s",
// 		conf.AWSRegion, conf.StorageServiceEndpoint, conf.BucketName)

// 	// Set Log Level
// 	logLvl, err := log.ParseLevel(conf.LogLevel)
// 	if err != nil {
// 		log.WithError(err).Error("unable to parse log level from config")

// 		return
// 	}

// 	log.SetOutput(os.Stdout)
// 	log.SetLevel(logLvl)

// 	// Create AWS configuration
// 	var awsConfig aws.Config

// 	if conf.StorageServiceEndpoint != "" {
// 		awsConfig = aws.Config{
// 			Endpoint:         aws.String(conf.StorageServiceEndpoint),
// 			Region:           aws.String(conf.AWSRegion),
// 			S3FOrcePathStyle: aws.Bool(true),
// 		}
// 	} else {
// 		awsConfig = aws.Config{
// 			Region:           aws.String(conf.AWSRegion),
// 			S3FOrcePathStyle: aws.Bool(true),
// 		}
// 	}

// 	// Create AWS session
// 	awsSession, err := session.NewSession(&awsConfig)
// 	if err != nil {
// 		log.Error("unable to make aws session", err)
// 		log.Printf("%+v\n", errors.WithStack(err))

// 		return
// 	}

// 	// Create S3 client
// 	s3Client := s3.New(awsSession)

// 	// Setup S3 bucket
// 	err = setupS3Bucket(conf, s3Client)
// 	if err != nil {
// 		logrus.Error("unable to setup the bucket: ", err)
// 		log.Printf("%+v\n", errors.WithStack(err))

// 		return
// 	}

// 	// TODO: Migrate object to S3 bucket
// }

// func setupS3Bucket(conf *config.MigrationTool, client *s3.S3) error {
// 	_, err := client.GetBucketAcl(&s3.GetBucketAclInput{
// 		Bucket: aws.String(conf.BucketName),
// 	})
// 	if err != nil {
// 		if err.(awserr.Error).Code() == s3.ErrCodeNoSuchBucket {
// 			err = createS3Bucket(conf, client)
// 			if err != nil {
// 				return errors.Wrap(err, "unable to create s3 bucket")
// 			}
// 		}
// 	}

// 	return nil
// }

// func createS3Bucket(conf *config.MigrationTool, client *s3.S3) error {
// 	input := &s3.CreateBucketInput{
// 		Bucket: &conf.BucketName,
// 		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
// 			LocationConstraint: aws.String(conf.AWSRegion),
// 		},
// 	}

// 	_, err := client.CreateBucket(input)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {
// 			switch awsErr.Code() {
// 			case s3.ErrCodeBucketAlreadyOwnedByYou:
// 				return nil
// 			case s3.ErrCodeBucketAlreadyExists:
// 				return nil
// 			default:
// 				return errors.Wrap(err, "unable to create bucket")
// 			}
// 		}
// 	}

// 	return nil
// }

// func checkBucketCreated(conf *config.MigrationTool, client *s3.S3) error {
// 	input := &s3.ListBucketsInput{}

// 	for {
// 		result, err := client.ListBuckets(input)
// 		if err != nil {
// 			return errors.Wrap(err, "unable to list bucket")
// 		}

// 		for _, bucket := range result.Buckets {
// 			if aws.StringValue(bucket.Name) == conf.BucketName {
// 				logrus.Infof("bucket %s created \n", conf.BucketName)

// 				return nil
// 			}
// 		}
// 	}
// }
