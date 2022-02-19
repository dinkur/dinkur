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

	dinkurapiv1 "github.com/dinkur/dinkur/api/dinkurapi/v1"
	"github.com/dinkur/dinkur/pkg/conv"
	"github.com/dinkur/dinkur/pkg/fromgrpc"
	"github.com/dinkur/dinkur/pkg/togrpc"
)

func (d *daemon) StreamAlert(req *dinkurapiv1.StreamAlertRequest, stream dinkurapiv1.Alerter_StreamAlertServer) error {
	if err := d.assertConnected(); err != nil {
		return convError(err)
	}
	if req == nil {
		return convError(ErrRequestIsNil)
	}
	ctx := stream.Context()
	ch, err := d.client.StreamAlert(ctx)
	if err != nil {
		return err
	}
	done := ctx.Done()
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(&dinkurapiv1.StreamAlertResponse{
				Alert: togrpc.Alert(ev.Alert),
				Event: togrpc.Event(ev.Event),
			}); err != nil {
				return err
			}
		case <-done:
			return nil
		}
	}
}

func (d *daemon) GetAlertList(ctx context.Context, req *dinkurapiv1.GetAlertListRequest) (*dinkurapiv1.GetAlertListResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	alerts, err := d.client.GetAlertList(ctx)
	if err != nil {
		return nil, err
	}
	return &dinkurapiv1.GetAlertListResponse{
		Alerts: togrpc.AlertSlice(alerts),
	}, nil
}

func (d *daemon) DeleteAlert(ctx context.Context, req *dinkurapiv1.DeleteAlertRequest) (*dinkurapiv1.DeleteAlertResponse, error) {
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
	deleted, err := d.client.DeleteAlert(ctx, id)
	if err != nil {
		return nil, convError(err)
	}
	return &dinkurapiv1.DeleteAlertResponse{
		DeletedAlert: togrpc.Alert(deleted),
	}, nil
}

func (d *daemon) DeleteAlertType(ctx context.Context, req *dinkurapiv1.DeleteAlertTypeRequest) (*dinkurapiv1.DeleteAlertTypeResponse, error) {
	if err := d.assertConnected(); err != nil {
		return nil, convError(err)
	}
	if req == nil {
		return nil, convError(ErrRequestIsNil)
	}
	deleted, err := d.client.DeleteAlertByType(ctx, fromgrpc.AlertType(req.Type))
	if err != nil {
		return nil, convError(err)
	}
	return &dinkurapiv1.DeleteAlertTypeResponse{
		DeletedAlert: togrpc.Alert(deleted),
	}, nil
}
