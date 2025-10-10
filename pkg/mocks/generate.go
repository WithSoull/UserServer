package mocks

//go:generate ../../bin/minimock -i github.com/WithSoull/UserServer/internal/service.UserService -o . -s "_minimock.go"
//go:generate ../../bin/minimock -i github.com/WithSoull/UserServer/internal/repository.UserRepository -o . -s "_minimock.go"
//go:generate ../../bin/minimock -i github.com/WithSoull/platform_common/pkg/client/db.TxManager -o . -s "_minimock.go"
