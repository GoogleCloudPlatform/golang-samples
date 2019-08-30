# gcputil/env

## Import


```shell
import "github.com/mchmarny/gcputil/env"
```

## Usage

Parse string environment variable

```shell
name := env.MustGetEnvVar("ENV_VAR_NAME", "default-value")
```

Parse int environment variable

```shell
name := env.MustGetIntEnvVar("HTTP_PORT", 8080)
```