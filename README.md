i2p-control
===========

Terminal interface to monitor and manage I2P router service. Basically, an
terminal i2pcontrol client.

        -host default:"127.0.0.1"
        -port default:"7657"
        -path default:"jsonrpc"
        -password default:"itoopie"
        -method default:"echo"
        -block default:false

Installation with go get

        go get -u github.com/eyedeekay/i2p-control

The methods that have been implemented are

        echo              : i2pcontrol:Echo
        restart           : i2pcontrol:Restart
        graceful-restart  : i2pcontrol:RestartGraceful
        shutdown          : i2pcontrol:Shutdown
        graceful-shutdown : i2pcontrol:ShutdownGraceful
        update            : i2pcontrol:Update
        find-update       : i2pcontrol:FindUpdate

So, for instance, to initiate a graceful shutdown and block until the router is
shut down, use the command:

        i2p-control -block -method=graceful-shutdown

