package model

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

type ClientEvent struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Version    string `json:"v"  binding:"required"`
	AppName    string `json:"an" binding:"required"`
	DataSource string `json:"ds" binding:"required"`
	AppVersion string `json:"av" binding:"required"`
	Type       string `json:"t"  binding:"required"`
	Nonce      string `json:"z"  binding:"required"`
	Sign       string `json:"sgn" binding:"required"`
	ClientID   string `json:"cid" binding:"required"`
	UserID     string `json:"uid"`

	Event      string            `json:"en" binding:"required"`
	EventValue map[string]string `json:"ev"`

	Timestamp time.Time
}

func (c *ClientEvent) Reset() {
	c.ID = ""
	c.Event = ""
	c.EventValue = nil
}

var ClientLogIndex mgo.Index = mgo.Index{
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
}
