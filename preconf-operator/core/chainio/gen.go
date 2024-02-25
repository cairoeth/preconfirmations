package chainio

//go:generate mockgen -destination=./mocks/avs_subscriber.go -package=mocks github.com/cairoeth/preconfirmations-avs/preconf-operator/core/chainio AvsSubscriberer
//go:generate mockgen -destination=./mocks/avs_writer.go -package=mocks github.com/cairoeth/preconfirmations-avs/preconf-operator/core/chainio AvsWriterer
//go:generate mockgen -destination=./mocks/avs_reader.go -package=mocks github.com/cairoeth/preconfirmations-avs/preconf-operator/core/chainio AvsReaderer
