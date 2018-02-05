package adapter

//校验AdapterConsole符合Adapter接口
var _ Adapter = &AdapterConsole{}
type AdapterConsole struct{}

func (ac *AdapterConsole) Write(b []byte) error {
    return nil
}

func (ac *AdapterConsole) Close() {

}