package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

type Audio struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name" json:"name"`
	FilePath   string             `bson:"path" json:"path"`
	StorageUrl string             `bson:"storage_url" json:"storage_url"`
}

type Job struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	AudioID         primitive.ObjectID `bson:"audio_id" json:"audio_id"`
	STTProcessed    bool               `bson:"initial_stt_done" json:"initial_stt_done"`
	STTtranscript   string             `bson:"stt_transcript" json:"stt_transcript"`
	HumanProcessed  bool               `bson:"human_processed" json:"human_processed"`
	HumanTranscript string             `bson:"human_transcript" json:"human_transcript"`
	WorkerID        primitive.ObjectID `bson:"worker_id" json:"worker_id"`
}
