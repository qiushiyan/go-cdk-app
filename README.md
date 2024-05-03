Create the app

```
mkdir go-aws
cd go-aws
cdk init app --language go
```

Prepare the environment

```
cdk bootstrap
```

Edit `go-aws.go` to refine the stack.

Then

```
cdk diff
cdk deploy
```
