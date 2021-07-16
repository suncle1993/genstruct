# genstruct

Golang struct generator from mysql schema

[![asciicast](https://asciinema.org/a/X5sk7TqrTTjF8AhN764K0Fc6m.svg)](https://asciinema.org/a/X5sk7TqrTTjF8AhN764K0Fc6m)

## 命令行版本

安装：

```
go get github.com/suncle1993/genstruct
```

使用方法：

```
genstruct -h 127.0.0.1 -u root -P 123456 -p 3306
```

* `-h` default `localhost`
* `-u` default `root`
* `-p` default `3306`

## 线上版本

- https://genstruct.suncle.me/


## 接口版本

```bash
curl --location --request GET 'https://genstructapi.herokuapp.com/api/struct/generate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "tags": ["db", "json"],
    "table": "create table user_mine_info( id bigint(20) NOT NULL AUTO_INCREMENT, uid bigint(20) NOT NULL DEFAULT '\''0'\'' COMMENT '\''用户uid'\'', mined_cnt bigint(20) NOT NULL COMMENT '\''剩余挖矿次数'\'', un_exchange_diamond bigint(20) NOT NULL COMMENT '\''未兑换为挖矿次数的钻石'\'', created_at bigint(20) NOT NULL COMMENT '\''创建时间'\'', updated_at bigint(20) NOT NULL COMMENT '\''更新时间'\'', PRIMARY KEY (id), UNIQUE KEY uk_uid (uid) USING BTREE) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '\''用户挖矿剩余次数记录'\'';"
}'
```

## 示例

建表数据

```mysql
CREATE TABLE user_mine_info(
  id bigint(20) NOT NULL AUTO_INCREMENT, 
  UID bigint(20) NOT NULL DEFAULT '0' COMMENT '用户uid', 
  mined_cnt bigint(20) NOT NULL COMMENT '剩余挖矿次数', 
  un_exchange_diamond bigint(20) NOT NULL COMMENT '未兑换为挖矿次数的钻石', 
  created_at bigint(20) NOT NULL COMMENT '创建时间', 
  updated_at bigint(20) NOT NULL COMMENT '更新时间', 
  PRIMARY KEY (id), 
  UNIQUE KEY uk_uid (UID) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户挖矿剩余次数记录';

```

生成的模型：

```go
package user_mine_info

// UserMineInfo 用户挖矿剩余次数记录
type UserMineInfo struct {
	ID                int64 `db:"id" json:"id" `
	UID               int64 `db:"uid" json:"uid" `                                 // 用户uid
	MinedCnt          int64 `db:"mined_cnt" json:"mined_cnt" `                     // 剩余挖矿次数
	UnExchangeDiamond int64 `db:"un_exchange_diamond" json:"un_exchange_diamond" ` // 未兑换为挖矿次数的钻石
	CreatedAt         int64 `db:"created_at" json:"created_at" `                   // 创建时间
	UpdatedAt         int64 `db:"updated_at" json:"updated_at" `                   // 更新时间
}

// TableName ...
func (u *UserMineInfo) TableName() string {
	return "user_mine_info" // TODO: 如果分表需要修改
}

// PK ...
func (u *UserMineInfo) PK() string {
	return "id"
}

// Schema ...
func (u *UserMineInfo) Schema() string {
	return `(
  id bigint NOT NULL AUTO_INCREMENT,
  uid bigint NOT NULL DEFAULT '0' COMMENT '用户uid',
  mined_cnt bigint NOT NULL COMMENT '剩余挖矿次数',
  un_exchange_diamond bigint NOT NULL COMMENT '未兑换为挖矿次数的钻石',
  created_at bigint NOT NULL COMMENT '创建时间',
  updated_at bigint NOT NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_uid (uid) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户挖矿剩余次数记录'`
}
```

