package dbclient

import (
	// "fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestParamsPairs(t *testing.T) {
	Convey("ParamsPairs should return a mapped params", t, func() {
		ret := ParamsPairs("k1", "v1", "k2", "v2", "k3")
		So(len(ret), ShouldEqual, 2)
	})
}

func TestBuildInsert(t *testing.T) {
	Convey("Insert builder should return sql and vals", t, func() {
		table := "table"
		pairs := ParamsPairs("k1", "v1", "k2", "v2", "k3")
		sql, vals := BuildInsert(table, pairs)
		sql = strings.ToLower(sql)
		isSql := sql == "insert into table(k1,k2) values(?,?)" || sql == "insert into table(k2,k1) values(?,?)"
		So(isSql, ShouldEqual, true)
		So(len(vals), ShouldEqual, 2)
	})
}

func TestBuildUpdate(t *testing.T) {
	Convey("Update builder should return sql and vals or panic", t, func() {
		Convey("Update builder should return sql and vals", func() {
			table := "table"
			pairs := ParamsPairs("k1", "v1", "k2", "v2", "k3")
			pairsConds := ParamsPairs("k5", "v5", "k6", "v6")
			sql, vals := BuildUpdate(table, pairs, pairsConds)
			sql = strings.ToLower(sql)
			isSql := sql == "update table set k1=?,k2=? where k5=? and k6=?" || sql == "update table set k2=?,k1=? where k5=? and k6=?"
			So(isSql, ShouldEqual, true)
			So(len(vals), ShouldEqual, 4)
		})

		Convey("Update builder should panic", func() {
			table := "table"
			pairs := ParamsPairs("k1", "v1", "k2", "v2", "k3")
			pairsConds := ParamsPairs()
			So(func() { BuildUpdate(table, pairs, pairsConds) }, ShouldPanic)
		})
	})
}

func BenchmarkParamsPairs(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParamsPairs("k1", "v1", "k2", "v2", "k3")
	}
}

func BenchmarkBuildInsert(b *testing.B) {
	table := "table"
	pairs := ParamsPairs("k1", "v1", "k2", "v2", "k3", "v3", "k4", "v4", "k5", "v5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildInsert(table, pairs)
	}
}

func BenchmarkBuildUpdate(b *testing.B) {
	table := "table"
	pairs := ParamsPairs("k1", "v1", "k2", "v2", "k3", "v3", "k4", "v4", "k5", "v5")
	pairsConds := ParamsPairs("k5", "v5", "k6", "v6")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildUpdate(table, pairs, pairsConds)
	}
}
