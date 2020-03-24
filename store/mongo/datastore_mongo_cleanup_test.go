// Copyright 2019 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/mendersoftware/deviceauth/model"
	"github.com/mendersoftware/deviceauth/store"
	"github.com/mendersoftware/go-lib-micro/identity"
	ctxstore "github.com/mendersoftware/go-lib-micro/store"
)

func TestGetDevicesBeingDecommissioned(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestGetDevicesBeingDecommissioned in short mode.")
	}

	testCases := []struct {
		inDevices  bson.A
		outDevices []model.Device
		tenant     string
	}{
		{
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000000",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
			outDevices: []model.Device{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
		},
		{
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000000",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
			outDevices: []model.Device{
				model.Device{
					Id: "00000000-0000-4000-8000-000000000001",
				},
			},
			tenant: tenant,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			testDbName := DbName
			if tc.tenant != "" {
				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
			}

			ctx := context.Background()
			if tc.tenant != "" {
				ctx = identity.WithContext(ctx, &identity.Identity{
					Tenant: tc.tenant,
				})
			}

			db := getDb(ctx)

			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
			_, err := coll.InsertMany(ctx, tc.inDevices)
			assert.NoError(t, err)

			brokenDevices, err := db.GetDevicesBeingDecommissioned(testDbName)
			assert.NoError(t, err)
			assert.Equal(t, tc.outDevices[0].Id, brokenDevices[0].Id)
		})
	}
}

func TestDeleteDevicesBeingDecommissioned(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestDeleteDevicesBeingDecommissioned in short mode.")
	}

	testCases := []struct {
		inDevices bson.A
		tenant    string
	}{
		{
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
		},
		{
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
			tenant: tenant,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			t.Logf("devices: %v", tc.inDevices)
			testDbName := DbName
			if tc.tenant != "" {
				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
			}

			ctx := context.Background()
			if tc.tenant != "" {
				ctx = identity.WithContext(ctx, &identity.Identity{
					Tenant: tc.tenant,
				})
			}

			db := getDb(ctx)

			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
			_, err := coll.InsertMany(ctx, tc.inDevices)
			assert.NoError(t, err)

			err = db.DeleteDevicesBeingDecommissioned(testDbName)
			assert.NoError(t, err)

			dbDevs, err := db.GetDevices(ctx, 0, 5, store.DeviceFilter{})
			assert.NoError(t, err)
			for _, dbDev := range dbDevs {
				assert.Equal(t, false, dbDev.Decommissioning)
			}
		})
	}
}

func TestGetBrokenAuthSets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestGetBrokenAuthSets in short mode.")
	}

	testCases := []struct {
		inAuthSets     bson.A
		inDevices      bson.A
		outAuthSetsIds []string
		tenant         string
		err            string
	}{
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000003",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
			},
			outAuthSetsIds: []string{"00000000-0000-4000-8000-000000000002"},
			err:            "",
		},
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000003",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
			},
			tenant:         tenant,
			outAuthSetsIds: []string{"00000000-0000-4000-8000-000000000002"},
			err:            "",
		},
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000002",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
			tenant:         tenant,
			outAuthSetsIds: []string{"00000000-0000-4000-8000-000000000002"},
			err:            "",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			testDbName := DbName
			if tc.tenant != "" {
				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
			}

			ctx := context.Background()
			if tc.tenant != "" {
				ctx = identity.WithContext(ctx, &identity.Identity{
					Tenant: tc.tenant,
				})
			}

			db := getDb(ctx)

			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbAuthSetColl)
			_, err := coll.InsertMany(ctx, tc.inAuthSets)
			assert.NoError(t, err)

			coll = db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
			_, err = coll.InsertMany(ctx, tc.inDevices)
			assert.NoError(t, err)

			brokenAuthSetsIds, err := db.GetBrokenAuthSets(testDbName)
			if tc.err != "" {
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.outAuthSetsIds, brokenAuthSetsIds)
			}
		})
	}
}

func TestDeleteBrokenAuthSets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestDeleteBrokenAuthSets in short mode.")
	}

	testCases := []struct {
		inAuthSets bson.A
		inDevices  bson.A
		tenant     string
	}{
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000003",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
			},
		},
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000003",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
			},
			tenant: tenant,
		},
		{
			inAuthSets: bson.A{
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000001",
					DeviceId: "00000000-0000-4000-8000-000000000001",
					IdData:   "001",
					PubKey:   "001",
				},
				model.AuthSet{
					Id:       "00000000-0000-4000-8000-000000000002",
					DeviceId: "00000000-0000-4000-8000-000000000002",
					IdData:   "001",
					PubKey:   "002",
				},
			},
			inDevices: bson.A{
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000001",
					IdData:          "001",
					PubKey:          "001",
					Status:          model.DevStatusPending,
					Decommissioning: false,
				},
				model.Device{
					Id:              "00000000-0000-4000-8000-000000000002",
					IdData:          "002",
					PubKey:          "002",
					Status:          model.DevStatusPending,
					Decommissioning: true,
				},
			},
			tenant: tenant,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			testDbName := DbName
			if tc.tenant != "" {
				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
			}

			ctx := context.Background()
			if tc.tenant != "" {
				ctx = identity.WithContext(ctx, &identity.Identity{
					Tenant: tc.tenant,
				})
			}

			db := getDb(ctx)

			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbAuthSetColl)
			_, err := coll.InsertMany(ctx, tc.inAuthSets)
			assert.NoError(t, err)

			coll = db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
			_, err = coll.InsertMany(ctx, tc.inDevices)
			assert.NoError(t, err)

			err = db.DeleteBrokenAuthSets(testDbName)
			assert.NoError(t, err)

			brokenAuthSetsIds, err := db.GetBrokenAuthSets(testDbName)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(brokenAuthSetsIds))
		})
	}
}

//func TestGetBrokenTokens(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping TestGetBrokenTokens in short mode.")
//	}
//
//	testCases := []struct {
//		inTokens     bson.A
//		inDevices    bson.A
//		outTokensIds []string
//		tenant       string
//	}{
//		{
//			inTokens: bson.A{
//				jwt.Token{jwt.Claims{
//					ID: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//				}},
//				jwt.Token{jwt.Claims{
//					Id: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000002")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000003")),
//				}},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//			},
//			outTokensIds: []string{"00000000-0000-4000-8000-000000000002"},
//		},
//		{
//			inTokens: bson.A{
//				jwt.Token{jwt.Claims{
//					ID: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//				}},
//				jwt.Token{jwt.Claims{
//					Id: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000002")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000003")),
//				}},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//			},
//			tenant:       tenant,
//			outTokensIds: []string{"00000000-0000-4000-8000-000000000002"},
//		},
//		{
//			inTokens: bson.A{
//				jwt.Token{jwt.Claims{
//					ID: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000001")),
//				}},
//				jwt.Token{jwt.Claims{
//					Id: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000002")),
//					Subject: uuid.Must(uuid.FromString(
//						"00000000-0000-4000-8000-000000000002")),
//				}},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: true,
//				},
//			},
//			tenant:       tenant,
//			outTokensIds: []string{"00000000-0000-4000-8000-000000000002"},
//		},
//	}
//
//	for i, tc := range testCases {
//		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
//			testDbName := DbName
//			if tc.tenant != "" {
//				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
//			}
//
//			ctx := context.Background()
//			if tc.tenant != "" {
//				ctx = identity.WithContext(ctx, &identity.Identity{
//					Tenant: tc.tenant,
//				})
//			}
//
//			db := getDb(ctx)
//
//			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbTokensColl)
//			_, err := coll.InsertMany(ctx, tc.inTokens)
//			assert.NoError(t, err)
//
//			coll = db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
//			_, err = coll.InsertMany(ctx, tc.inDevices)
//			assert.NoError(t, err)
//
//			brokenTokensIds, err := db.GetBrokenTokens(testDbName)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.outTokensIds, brokenTokensIds)
//		})
//	}
//}

//func TestDeleteBrokenTokens(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping TestDeleteBrokenAuthSets in short mode.")
//	}
//
//	testCases := []struct {
//		inTokens  bson.A
//		inDevices bson.A
//		tenant    string
//	}{
//		{
//			inTokens: bson.A{
//				model.Token{
//					Id:        "001",
//					DevId:     "001",
//					AuthSetId: "001",
//					Token:     "foo",
//				},
//				model.Token{
//					Id:        "002",
//					DevId:     "003",
//					AuthSetId: "002",
//					Token:     "bar",
//				},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//			},
//		},
//		{
//			inTokens: bson.A{
//				model.Token{
//					Id:    "001",
//					DevId: "001",
//					Token: "foo",
//				},
//				model.Token{
//					Id:    "002",
//					DevId: "003",
//					Token: "bar",
//				},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//			},
//			tenant: tenant,
//		},
//		{
//			inTokens: bson.A{
//				model.Token{
//					Id:    "001",
//					DevId: "001",
//					Token: "foo",
//				},
//				model.Token{
//					Id:    "002",
//					DevId: "002",
//					Token: "bar",
//				},
//			},
//			inDevices: bson.A{
//				model.Device{
//					Id:              "001",
//					IdData:          "001",
//					PubKey:          "001",
//					Status:          model.DevStatusPending,
//					Decommissioning: false,
//				},
//				model.Device{
//					Id:              "002",
//					IdData:          "002",
//					PubKey:          "002",
//					Status:          model.DevStatusPending,
//					Decommissioning: true,
//				},
//			},
//			tenant: tenant,
//		},
//	}
//
//	for i, tc := range testCases {
//		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
//			testDbName := DbName
//			if tc.tenant != "" {
//				testDbName = ctxstore.DbNameForTenant(tc.tenant, DbName)
//			}
//
//			ctx := context.Background()
//			if tc.tenant != "" {
//				ctx = identity.WithContext(ctx, &identity.Identity{
//					Tenant: tc.tenant,
//				})
//			}
//
//			db := getDb(ctx)
//
//			coll := db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbTokensColl)
//			_, err := coll.InsertMany(ctx, tc.inTokens)
//			assert.NoError(t, err)
//
//			coll = db.client.Database(ctxstore.DbFromContext(ctx, DbName)).Collection(DbDevicesColl)
//			_, err = coll.InsertMany(ctx, tc.inDevices)
//			assert.NoError(t, err)
//
//			err = db.DeleteBrokenTokens(testDbName)
//			assert.NoError(t, err)
//
//			brokenTokensIds, err := db.GetBrokenTokens(testDbName)
//			assert.NoError(t, err)
//			assert.Equal(t, 0, len(brokenTokensIds))
//		})
//	}
//}
