// cfg 包使用 github.com/spf13/viper 提供配置能力，
// 目前按照优先级支持 env/file，先读取 env 的 BLOG_PRODUCTION
// true: 读取 cfg_prod.json
// false: 读取 cfg_dev.json
package cfg
