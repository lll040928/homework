package global

type Config struct {
	Host        string      `mapstructure:"host"`
	Port        string      `mapstructure:"port"`
	MysqlInfo   MySQLConfig `mapstructure:"mysql_info"`
	RedisInfo   RedisConfig `mapstructure:"redis_info"`
	ZapInfo     ZapConfig   `mapstructure:"zap_info"`
	CryptInfo   CryptConfig `mapstructure:"crypt_info"`
	TokenSecret string      `mapstructure:"token_secret"`
}
type ZapConfig struct {
	Path string `mapstructure:"path"`
}
type MySQLConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}

type CryptConfig struct {
	Key string `mapstructure:"key"`
	Iv  string `mapstructure:"iv"`
}
