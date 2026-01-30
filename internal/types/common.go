package types

const (
	MorphAPIURL  = "https://www.morphllm.com/api/warpgrep-chat"
	MorphCookies = "_gcl_aw=GCL.1769242305.Cj0KCQiA-NHLBhDSARIsAIhe9X36gjH12LKRxYZ1bqy4c8wATix3qqxauv9-7H-rewEfnKOW0npi5ZMaAjfHEALw_wcB; _gcl_gs=2.1.k1$i1769242301$u142736937; _ga=GA1.1.1601700642.1769242305; __client_uat=1769312595; __client_uat_B3JMRlZP=1769312595; __refresh_B3JMRlZP=mCRNoZazsM9bFPAAxEwv; clerk_active_context=sess_38jUhuq7Hkg8fiqhHMduDWfmycM:; _rdt_uuid=1769242305266.9f8d183d-a987-4e5b-9dee-336b47494719; _rdt_em=:3a02869661fc8f1da6409edc313bf251ebc4586f202b2b84a31646b48d1beca7,3a02869661fc8f1da6409edc313bf251ebc4586f202b2b84a31646b48d1beca7,320e81b4145c1933c25a4ee8275675397a23cb7594b9ac52006f470805bbbe42,3eb2031bc3c935a9a6f57367e9d851317a41280a8ac0693930942a8ee631e447,a03bc4dd2a023bebb17e8f64b22f9803685055600064656f8f85675a4dd3622d; __session=eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yc2NFSEpuWHRhREZVVXhVQ1habldocTVYS0MiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJodHRwczovL3d3dy5tb3JwaGxsbS5jb20iLCJleHAiOjE3NjkzMjU5MzIsImZ2YSI6WzIyMSwtMV0sImlhdCI6MTc2OTMyNTg3MiwiaXNzIjoiaHR0cHM6Ly9jbGVyay5tb3JwaGxsbS5jb20iLCJuYmYiOjE3NjkzMjU4NjIsInNpZCI6InNlc3NfMzhqVWh1cTdIa2c4ZmlxaEhNZHVEV2ZteWNNIiwic3RzIjoiYWN0aXZlIiwic3ViIjoidXNlcl8zOGpVaHRFRTBTTzNaaTFZaVhFZmI5UjQ3RVoifQ.HhJ-kSkAOlU7xofFZOjHEjH_RbzkEtJwdjxTJ_C09_VgLCOv1dd7OOKtMt5sRXpIl80aF4FrWt8_XaBd4bx63JDzZUC1len0b1CIubu1n6t6vNRts_0hHsaxVo6NTxiElJzZkBZfoZIdnUC5zsEPDldct3wW9C0Jb77Em_uDWlDvs1D-5UF0Hol75v8wj-dnys01IVUqN0svt-QlJ0mKTwbCMeFh4mh1UbE8rMKassgblVJpGfNWWr3pzscuS3yxRbvq9URrr-HoweybWlq57SaJJMxpwopGdnRx2jrKrw_IcJQlQ_Ug_8HxF69i3Upu_rVFIdbIx2hsV0o7LfHzGg; __session_B3JMRlZP=eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yc2NFSEpuWHRhREZVVXhVQ1habldocTVYS0MiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJodHRwczovL3d3dy5tb3JwaGxsbS5jb20iLCJleHAiOjE3NjkzMjU5MzIsImZ2YSI6WzIyMSwtMV0sImlhdCI6MTc2OTMyNTg3MiwiaXNzIjoiaHR0cHM6Ly9jbGVyay5tb3JwaGxsbS5jb20iLCJuYmYiOjE3NjkzMjU4NjIsInNpZCI6InNlc3NfMzhqVWh1cTdIa2c4ZmlxaEhNZHVEV2ZteWNNIiwic3RzIjoiYWN0aXZlIiwic3ViIjoidXNlcl8zOGpVaHRFRTBTTzNaaTFZaVhFZmI5UjQ3RVoifQ.HhJ-kSkAOlU7xofFZOjHEjH_RbzkEtJwdjxTJ_C09_VgLCOv1dd7OOKtMt5sRXpIl80aF4FrWt8_XaBd4bx63JDzZUC1len0b1CIubu1n6t6vNRts_0hHsaxVo6NTxiElJzZkBZfoZIdnUC5zsEPDldct3wW9C0Jb77Em_uDWlDvs1D-5UF0Hol75v8wj-dnys01IVUqN0svt-QlJ0mKTwbCMeFh4mh1UbE8rMKassgblVJpGfNWWr3pzscuS3yxRbvq9URrr-HoweybWlq57SaJJMxpwopGdnRx2jrKrw_IcJQlQ_Ug_8HxF69i3Upu_rVFIdbIx2hsV0o7LfHzGg; _ga_5ET01XBKB1=GS2.1.s1769324224$o6$g1$t1769325875$j59$l0$h0; ph_phc_i9YGegL9gG85W32ArqnIiBECNTAzUYlFFK8B0Odbhk8_posthog=%7B%22%24device_id%22%3A%22019bef0f-25f6-7e51-b718-8db6665d4470%22%2C%22distinct_id%22%3A%22user_38jUhtEE0SO3Zi1YiXEfb9R47EZ%22%2C%22%24sesid%22%3A%5B1769325875736%2C%22019bf3f1-202a-7fd6-8dbb-6cbea73ffa8f%22%2C1769324224553%5D%2C%22%24epp%22%3Atrue%2C%22%24initial_person_info%22%3A%7B%22r%22%3A%22https%3A%2F%2Fwww.google.com%2F%22%2C%22u%22%3A%22https%3A%2F%2Fwww.morphllm.com%2Fpricing%3Fgad_source%3D1%26gad_campaignid%3D23473818447%26gbraid%3D0AAAAA_wYB2fufbMARVR2fJKKHNzc8Ovu-%26gclid%3DCj0KCQiA-NHLBhDSARIsAIhe9X36gjH12LKRxYZ1bqy4c8wATix3qqxauv9-7H-rewEfnKOW0npi5ZMaAjfHEALw_wcB%22%7D%7D"
	LogDir       = "./logs"
)

var MorphHeaders = map[string]string{
	"accept":             "*/*",
	"accept-language":    "zh-CN,zh;q=0.9",
	"cache-control":      "no-cache",
	"content-type":       "application/json",
	"origin":             "https://www.morphllm.com",
	"pragma":             "no-cache",
	"priority":           "u=1, i",
	"referer":            "https://www.morphllm.com/playground/na/warpgrep?repo=tiangolo%2Ffastapi",
	"sec-ch-ua":          `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`,
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": `"macOS"`,
	"sec-fetch-dest":     "empty",
	"sec-fetch-mode":     "cors",
	"sec-fetch-site":     "same-origin",
	"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
	"cookie":             MorphCookies,
}

var DebugMode = true

type ParsedToolCall struct {
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}
