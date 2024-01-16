package rtclib

import (
	"context"

	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pingostack/neon/pkg/eventemitter"
	"github.com/pingostack/neon/pkg/logger"
	"github.com/pingostack/neon/protocols/rtclib/transport"
	"github.com/pkg/errors"
)

type LocalStream struct {
	*transport.Transport
	ctx          context.Context
	cancel       context.CancelFunc
	logger       logger.Logger
	eventemitter eventemitter.EventEmitter
}

func NewLocalStream(transport *transport.Transport) (*LocalStream, error) {
	c := &LocalStream{
		Transport:    transport,
		logger:       transport.Logger(),
		eventemitter: eventemitter.NewEventEmitter(transport.Context(), defaultEventEmitterLength, transport.Logger()),
	}

	c.ctx, c.cancel = context.WithCancel(transport.Context())

	if err := c.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid local stream")
	}

	return c, nil
}

func (c *LocalStream) validate() error {
	if c.Transport == nil {
		return errors.New("transport not set")
	}

	if c.ctx == nil {
		ctx, cancel := context.WithCancel(context.Background())
		c.ctx = ctx
		c.cancel = cancel
	}

	if c.logger == nil {
		c.logger = logger.DefaultLogger
	}

	return nil
}

func (c *LocalStream) SetupTracks(videoTrack format.Format, audioTrack format.Format) ([]*TrackLocl, error) {
	var tracks []*TrackLocl

	for _, forma := range []format.Format{videoTrack, audioTrack} {
		if forma != nil {
			track, err := NewTrackLocl(forma, c.Transport.AddTrack)
			if err != nil {
				return nil, err
			}

			tracks = append(tracks, track)
		}
	}

	return tracks, nil
}
