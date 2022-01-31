// Dinkur the task time tracking utility.
// <https://github.com/dinkur/dinkur>
//
// Copyright (C) 2021 Kalle Fagerberg
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

package dinkurdb

import (
	"context"
	"fmt"

	"github.com/dinkur/dinkur/pkg/dinkur"
	"gopkg.in/typ.v1"
	"gorm.io/gorm"
)

func (c *client) StreamAlert(ctx context.Context) (<-chan dinkur.StreamedAlert, error) {
	if err := c.assertConnected(); err != nil {
		return nil, err
	}
	ch := make(chan dinkur.StreamedAlert)
	go func() {
		done := ctx.Done()
		dbAlertChan := c.alertObs.Sub()
		defer close(ch)
		defer func() {
			if err := c.alertObs.Unsub(dbAlertChan); err != nil {
				log.Warn().WithError(err).Message("Failed to unsub alert.")
			}
		}()
		for {
			select {
			case ev, ok := <-dbAlertChan:
				if !ok {
					return
				}
				alert, err := convAlert(ev.dbAlert)
				if err != nil {
					log.Warn().
						WithError(err).
						WithUint("alertId", ev.dbAlert.ID).
						Message("Invalid alert event.")
					continue
				}
				ch <- dinkur.StreamedAlert{
					Alert: alert,
					Event: ev.event,
				}
			case <-done:
				return
			}
		}
	}()
	return ch, nil
}

func (c *client) GetAlertList(ctx context.Context) ([]dinkur.Alert, error) {
	if err := c.assertConnected(); err != nil {
		return nil, err
	}
	dbAlerts, err := c.listDBAlertsAtom()
	if err != nil {
		return nil, err
	}
	alerts, err := typ.MapErr(dbAlerts, convAlert)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (c *client) listDBAlertsAtom() ([]Alert, error) {
	if err := c.assertConnected(); err != nil {
		return nil, err
	}
	var dbAlerts []Alert
	if err := c.dbAlertPreloaded().Find(&dbAlerts).Error; err != nil {
		return nil, err
	}
	return dbAlerts, nil
}

func (c *client) DeleteAlert(ctx context.Context, id uint) (dinkur.Alert, error) {
	if err := c.assertConnected(); err != nil {
		return nil, err
	}
	dbAlert, err := c.withContext(ctx).deleteDBAlertAtom(id)
	if err != nil {
		return nil, err
	}
	c.alertObs.PubWait(alertEvent{
		dbAlert: dbAlert,
		event:   dinkur.EventDeleted,
	})
	return nil, nil
}

func (c *client) deleteDBAlertAtom(id uint) (Alert, error) {
	var dbAlert Alert
	err := c.transaction(func(tx *client) (tranErr error) {
		dbAlert, tranErr = tx.deleteDBAlertNoTran(id)
		return
	})
	return dbAlert, err
}

func (c *client) deleteDBAlertNoTran(id uint) (Alert, error) {
	dbAlert, err := c.getDBAlertAtom(id)
	if err != nil {
		return Alert{}, fmt.Errorf("get alert to delete: %w", err)
	}
	if err := c.db.Delete(&Entry{}, id).Error; err != nil {
		return Alert{}, fmt.Errorf("delete alert: %w", err)
	}
	return dbAlert, nil
}

func (c *client) getDBAlertAtom(id uint) (Alert, error) {
	if err := c.assertConnected(); err != nil {
		return Alert{}, err
	}
	var dbAlert Alert
	err := c.dbAlertPreloaded().First(&dbAlert, id).Error
	if err != nil {
		return Alert{}, err
	}
	return dbAlert, nil
}

func (c *client) dbAlertPreloaded() *gorm.DB {
	return c.db.Model(&Alert{}).
		Preload(alertColumnPlainMessage).
		Preload(alertColumnAFK)
}
