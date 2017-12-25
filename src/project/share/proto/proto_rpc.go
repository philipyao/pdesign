package proto

type GameHelloArg struct {
    A, B int
}
type GameHelloRep struct {
    C int
}

type FetchConfigArg struct {
    Namespace       string
    Keys            []string
}
type FetchConfigRes struct {
    Errmsg          string
    Confs           []*ConfigEntry
}
type ConfigEntry struct {
    Key             string
    Value           string
}