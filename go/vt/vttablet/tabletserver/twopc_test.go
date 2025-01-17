/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tabletserver

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/test/utils"
	querypb "vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/tx"
)

func TestReadAllRedo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Reuse code from tx_executor_test.
	_, tsv, db := newTestTxExecutor(t, ctx)
	defer db.Close()
	defer tsv.StopService()
	tpc := tsv.te.twoPC

	conn, err := tsv.qe.conns.Get(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Recycle()

	db.AddQuery(tpc.readAllRedo, &sqltypes.Result{})
	prepared, failed, err := tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var want []*tx.PreparedTx
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt01"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt01"),
		}, {
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt02"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt01"),
		}, {
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt02"),
		}, {
			sqltypes.NewVarBinary("dtid1"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt11"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid1",
		Queries: []string{"stmt11"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt01"),
		}, {
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt02"),
		}, {
			sqltypes.NewVarBinary("dtid1"),
			sqltypes.NewVarBinary("Failed"),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt11"),
		}, {
			sqltypes.NewVarBinary("dtid2"),
			sqltypes.NewVarBinary("Failed"),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt21"),
		}, {
			sqltypes.NewVarBinary("dtid2"),
			sqltypes.NewVarBinary("Failed"),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt22"),
		}, {
			sqltypes.NewVarBinary("dtid3"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("stmt31"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid3",
		Queries: []string{"stmt31"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	wantFailed := []*tx.PreparedTx{{
		Dtid:    "dtid1",
		Queries: []string{"stmt11"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid2",
		Queries: []string{"stmt21", "stmt22"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(failed, wantFailed) {
		t.Errorf("ReadAllRedo failed): %s, want %s", jsonStr(failed), jsonStr(wantFailed))
	}
}

func TestReadAllTransactions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, tsv, db := newTestTxExecutor(t, ctx)
	defer db.Close()
	defer tsv.StopService()
	tpc := tsv.te.twoPC

	conn, err := tsv.qe.conns.Get(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Recycle()

	db.AddQuery(tpc.readAllTransactions, &sqltypes.Result{})
	distributed, err := tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var want []*tx.DistributedTx
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions: %s, want %s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks01"),
			sqltypes.NewVarBinary("shard01"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []querypb.Target{{
			Keyspace: "ks01",
			Shard:    "shard01",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks01"),
			sqltypes.NewVarBinary("shard01"),
		}, {
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks02"),
			sqltypes.NewVarBinary("shard02"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []querypb.Target{{
			Keyspace: "ks01",
			Shard:    "shard01",
		}, {
			Keyspace: "ks02",
			Shard:    "shard02",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &sqltypes.Result{
		Fields: []*querypb.Field{
			{Type: sqltypes.VarChar},
			{Type: sqltypes.Int64},
			{Type: sqltypes.Int64},
			{Type: sqltypes.VarChar},
			{Type: sqltypes.VarChar},
		},
		Rows: [][]sqltypes.Value{{
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks01"),
			sqltypes.NewVarBinary("shard01"),
		}, {
			sqltypes.NewVarBinary("dtid0"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks02"),
			sqltypes.NewVarBinary("shard02"),
		}, {
			sqltypes.NewVarBinary("dtid1"),
			sqltypes.NewInt64(RedoStatePrepared),
			sqltypes.NewVarBinary("1"),
			sqltypes.NewVarBinary("ks11"),
			sqltypes.NewVarBinary("shard11"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []querypb.Target{{
			Keyspace: "ks01",
			Shard:    "shard01",
		}, {
			Keyspace: "ks02",
			Shard:    "shard02",
		}},
	}, {
		Dtid:    "dtid1",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []querypb.Target{{
			Keyspace: "ks11",
			Shard:    "shard11",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}
}

func jsonStr(v any) string {
	out, _ := json.Marshal(v)
	return string(out)
}

// TestUnresolvedTransactions tests the retrieval of unresolved transactions from the database and
// providing the output in proto format.
func TestUnresolvedTransactions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, tsv, db := newTestTxExecutor(t, ctx)
	defer db.Close()
	defer tsv.StopService()

	conn, err := tsv.qe.conns.Get(ctx, nil)
	require.NoError(t, err)
	defer conn.Recycle()

	tcases := []struct {
		name         string
		unresolvedTx *sqltypes.Result
		expectedTx   []*querypb.TransactionMetadata
	}{{
		name:         "no unresolved transactions",
		unresolvedTx: &sqltypes.Result{},
	}, {
		name: "one unresolved transaction",
		unresolvedTx: sqltypes.MakeTestResult(
			sqltypes.MakeTestFields("dtid|state|keyspace|shard",
				"VARBINARY|INT64|VARCHAR|VARCHAR"),
			"dtid0|1|ks01|shard01",
			"dtid0|1|ks01|shard02"),
		expectedTx: []*querypb.TransactionMetadata{{
			Dtid:  "dtid0",
			State: querypb.TransactionState_PREPARE,
			Participants: []*querypb.Target{
				{Keyspace: "ks01", Shard: "shard01"},
				{Keyspace: "ks01", Shard: "shard02"},
			}}},
	}, {
		name: "two unresolved transaction",
		unresolvedTx: sqltypes.MakeTestResult(
			sqltypes.MakeTestFields("dtid|state|keyspace|shard",
				"VARBINARY|INT64|VARCHAR|VARCHAR"),
			"dtid0|2|ks01|shard01",
			"dtid0|2|ks01|shard02",
			"dtid1|3|ks02|shard03",
			"dtid1|3|ks01|shard02"),
		expectedTx: []*querypb.TransactionMetadata{{
			Dtid:  "dtid0",
			State: querypb.TransactionState_COMMIT,
			Participants: []*querypb.Target{
				{Keyspace: "ks01", Shard: "shard01"},
				{Keyspace: "ks01", Shard: "shard02"},
			}}, {
			Dtid:  "dtid1",
			State: querypb.TransactionState_ROLLBACK,
			Participants: []*querypb.Target{
				{Keyspace: "ks02", Shard: "shard03"},
				{Keyspace: "ks01", Shard: "shard02"},
			}},
		},
	}}

	tpc := tsv.te.twoPC
	txQueryPattern := `.*time_created < 1000.*`
	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			db.AddQueryPattern(txQueryPattern, tcase.unresolvedTx)
			distributed, err := tpc.UnresolvedTransactions(ctx, time.UnixMicro(1))
			require.NoError(t, err)
			utils.MustMatch(t, tcase.expectedTx, distributed)
		})
	}
}
