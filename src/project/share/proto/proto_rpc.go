package proto

type GameHelloArg struct {
    A, B int
}
type GameHelloRep struct {
    C int
}

type ConfigWithNamespaceArg struct {
    Namespace       string
}
type ConfigWithNamespaceRep struct {
    Confs           []*Config
}
type Config struct {
    Key             string
    Value           string
}