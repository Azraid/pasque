start ./router.exe Router.1 ./config/system.json
start ./gate.exe MobileTrade.Gate ./config/system.json
start ./gate.exe B2.ApiGate ./config/system.json
start ./fakesrv.exe MobileTradeSrv.1 ./config/system.json ./scenario.json
