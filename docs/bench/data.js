window.BENCHMARK_DATA = {
  "lastUpdate": 1789399437155,
  "repoUrl": "https://github.com/moov-io/iso8583",
  "entries": {
    "moov-io/iso8583": [
      {
        "commit": {
          "author": {
            "name": "Adam Shannon",
            "username": "adamdecaf",
            "email": "adamkshannon@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "65dfa97dd20ca9d33d5e8f02354d40d117c75e83",
          "message": "Track marshal and unpack Go benchmarks in this repository. (#453)\n\nStore results in docs/bench labeled with the iso8583 commit that ran.",
          "timestamp": "2026-09-14T15:13:58Z",
          "url": "https://github.com/moov-io/iso8583/commit/65dfa97dd20ca9d33d5e8f02354d40d117c75e83"
        },
        "date": 1789399436465,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkMarshaling",
            "value": 19731,
            "unit": "ns/op\t   13680 B/op\t     275 allocs/op",
            "extra": "60682 times\n4 procs"
          },
          {
            "name": "BenchmarkMarshaling - ns/op",
            "value": 19731,
            "unit": "ns/op",
            "extra": "60682 times\n4 procs"
          },
          {
            "name": "BenchmarkMarshaling - B/op",
            "value": 13680,
            "unit": "B/op",
            "extra": "60682 times\n4 procs"
          },
          {
            "name": "BenchmarkMarshaling - allocs/op",
            "value": 275,
            "unit": "allocs/op",
            "extra": "60682 times\n4 procs"
          },
          {
            "name": "BenchmarkUnpacking",
            "value": 14174,
            "unit": "ns/op\t   14784 B/op\t     286 allocs/op",
            "extra": "82663 times\n4 procs"
          },
          {
            "name": "BenchmarkUnpacking - ns/op",
            "value": 14174,
            "unit": "ns/op",
            "extra": "82663 times\n4 procs"
          },
          {
            "name": "BenchmarkUnpacking - B/op",
            "value": 14784,
            "unit": "B/op",
            "extra": "82663 times\n4 procs"
          },
          {
            "name": "BenchmarkUnpacking - allocs/op",
            "value": 286,
            "unit": "allocs/op",
            "extra": "82663 times\n4 procs"
          }
        ]
      }
    ]
  }
}