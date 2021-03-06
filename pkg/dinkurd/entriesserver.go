// Dinkur the task time tracking utility.
// <https://github.com/dinkur/dinkur>
//
// SPDX-FileCopyrightText: 2021 Kalle Fagerberg
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package dinkurd

import (
	"context"
	"errors"

	dinkurapiv1 "github.com/dinkur/dinkur/api/dinkurapi/v1"
	"github.com/dinkur/dinkur/pkg/conv"
	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/fromgrpc"
	"github.com/dinkur/dinkur/pkg/togrpc"
)

func (d *daemon) Ping(ctx context.Context, req *dinkurapiv1.PingRequest) (*dinkurapiv1.PingResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	if err := d.client.Ping(ctx); err != nil {
		return nil, convError(err)
	}
	return &dinkurapiv1.PingResponse{}, nil
}

func (d *daemon) GetEntry(ctx context.Context, req *dinkurapiv1.GetEntryRequest) (*dinkurapiv1.GetEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	id, err := conv.Uint64ToUint(req.Id)
	if err != nil {
		return nil, convError(err)
	}
	entry, err := d.client.GetEntry(ctx, id)
	if err != nil {
		if errors.Is(err, dinkur.ErrNotFound) {
			return &dinkurapiv1.GetEntryResponse{}, nil
		}
		return nil, convError(err)
	}
	return &dinkurapiv1.GetEntryResponse{
		Entry: togrpc.EntryPtr(&entry),
	}, nil
}

func (d *daemon) GetActiveEntry(ctx context.Context, req *dinkurapiv1.GetActiveEntryRequest) (*dinkurapiv1.GetActiveEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	entry, err := d.client.GetActiveEntry(ctx)
	if err != nil {
		if errors.Is(err, dinkur.ErrNotFound) {
			return &dinkurapiv1.GetActiveEntryResponse{}, nil
		}
		return nil, convError(err)
	}
	return &dinkurapiv1.GetActiveEntryResponse{
		ActiveEntry: togrpc.EntryPtr(entry),
	}, nil
}

func (d *daemon) GetEntryList(ctx context.Context, req *dinkurapiv1.GetEntryListRequest) (*dinkurapiv1.GetEntryListResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	search := dinkur.SearchEntry{
		Start:              fromgrpc.TimePtr(req.Start),
		End:                fromgrpc.TimePtr(req.End),
		Shorthand:          fromgrpc.Shorthand(req.Shorthand),
		NameFuzzy:          req.NameFuzzy,
		NameHighlightStart: req.NameHighlightStart,
		NameHighlightEnd:   req.NameHighlightEnd,
	}
	var err error
	search.Limit, err = conv.Uint64ToUint(req.Limit)
	if err != nil {
		return nil, convError(err)
	}
	entries, err := d.client.GetEntryList(ctx, search)
	if err != nil {
		return nil, convError(err)
	}
	return &dinkurapiv1.GetEntryListResponse{
		Entries: togrpc.EntrySlice(entries),
	}, nil
}

func (d *daemon) CreateEntry(ctx context.Context, req *dinkurapiv1.CreateEntryRequest) (*dinkurapiv1.CreateEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	startAfterID, err := conv.Uint64ToUint(req.StartAfterIdOrZero)
	if err != nil {
		return nil, convError(err)
	}
	endBeforeID, err := conv.Uint64ToUint(req.EndBeforeIdOrZero)
	if err != nil {
		return nil, convError(err)
	}
	newEntry := dinkur.NewEntry{
		Name:               req.Name,
		Start:              fromgrpc.TimePtr(req.Start),
		End:                fromgrpc.TimePtr(req.End),
		StartAfterIDOrZero: startAfterID,
		EndBeforeIDOrZero:  endBeforeID,
		StartAfterLast:     req.StartAfterLast,
	}
	startedEntry, err := d.client.CreateEntry(ctx, newEntry)
	if err != nil {
		return nil, convError(err)
	}
	d.onEntryMutation(ctx)
	return &dinkurapiv1.CreateEntryResponse{
		PreviouslyActiveEntry: togrpc.EntryPtr(startedEntry.Stopped),
		CreatedEntry:          togrpc.EntryPtr(&startedEntry.Started),
	}, nil
}

func (d *daemon) UpdateEntry(ctx context.Context, req *dinkurapiv1.UpdateEntryRequest) (*dinkurapiv1.UpdateEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	id, err := conv.Uint64ToUint(req.IdOrZero)
	if err != nil {
		return nil, convError(err)
	}
	startAfterID, err := conv.Uint64ToUint(req.StartAfterIdOrZero)
	if err != nil {
		return nil, convError(err)
	}
	endBeforeID, err := conv.Uint64ToUint(req.EndBeforeIdOrZero)
	if err != nil {
		return nil, convError(err)
	}
	edit := dinkur.EditEntry{
		Name:               conv.ZeroAsNil(req.Name),
		Start:              fromgrpc.TimePtr(req.Start),
		StartFuzzy:         req.StartFuzzy,
		End:                fromgrpc.TimePtr(req.End),
		EndFuzzy:           req.EndFuzzy,
		IDOrZero:           id,
		AppendName:         req.AppendName,
		StartAfterIDOrZero: startAfterID,
		EndBeforeIDOrZero:  endBeforeID,
		StartAfterLast:     req.StartAfterLast,
	}
	update, err := d.client.UpdateEntry(ctx, edit)
	if err != nil {
		return nil, convError(err)
	}
	d.onEntryMutation(ctx)
	return &dinkurapiv1.UpdateEntryResponse{
		Before: togrpc.EntryPtr(&update.Before),
		After:  togrpc.EntryPtr(&update.After),
	}, nil
}

func (d *daemon) DeleteEntry(ctx context.Context, req *dinkurapiv1.DeleteEntryRequest) (*dinkurapiv1.DeleteEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	id, err := conv.Uint64ToUint(req.Id)
	if err != nil {
		return nil, convError(err)
	}
	deletedEntry, err := d.client.DeleteEntry(ctx, id)
	if err != nil {
		return nil, convError(err)
	}
	d.onEntryMutation(ctx)
	return &dinkurapiv1.DeleteEntryResponse{
		DeletedEntry: togrpc.EntryPtr(&deletedEntry),
	}, nil
}

func (d *daemon) StopActiveEntry(ctx context.Context, req *dinkurapiv1.StopActiveEntryRequest) (*dinkurapiv1.StopActiveEntryResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	stoppedEntry, err := d.client.StopActiveEntry(ctx, fromgrpc.TimeOrNow(req.End))
	if err != nil {
		return nil, convError(err)
	}
	d.onEntryMutation(ctx)
	return &dinkurapiv1.StopActiveEntryResponse{
		StoppedEntry: togrpc.EntryPtr(stoppedEntry),
	}, nil
}

func (d *daemon) StreamEntry(req *dinkurapiv1.StreamEntryRequest, stream dinkurapiv1.Entries_StreamEntryServer) error {
	if err := d.assertConnected(); err != nil {
		return convError(err)
	}
	if req == nil {
		return convError(ErrRequestIsNil)
	}
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	ch, err := d.client.StreamEntry(ctx)
	if err != nil {
		return convError(err)
	}
	for ev := range ch {
		if err := stream.Send(&dinkurapiv1.StreamEntryResponse{
			Entry: togrpc.EntryPtr(&ev.Entry),
			Event: togrpc.Event(ev.Event),
		}); err != nil {
			return convError(err)
		}
	}
	return nil
}
