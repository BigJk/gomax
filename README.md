gomax
====

gomax is a corewar optimizer written in go. The syntax and idea is taken from optiMAX.

## gomax cli

### create
```gomax create %output% %config% %warrior%```
creates the database that is needed for the optimization progress.

* %output% - Output file of the database. (eg. opt.db)
* %config% - Config file for the optimization. (eg. config.json)
* %warrior% - Warrior you want to optimize. (eg. warrior.red)

**example:** ```gomax create opt.db config.json warrior.red```

### run
```gomax run %database% %threads%```
starts the optimization progress with a set amount of threads.

* %database% - Database file. (eg. opt.db)
* %threads% - Count of threads. Max should be your CPU core count. (eg. 6)

**example:** ```gomax run opt.db 6```

### config
```gomax config```
creates a sample config.

**example:** ```gomax config```

## gomax webpanel

**visit the local webpanel under:** ```http://127.0.0.1```

## gomax uses

#### Back-end

* [BoltDB](https://github.com/boltdb/bolt)
* [exMars](http://corewar.co.uk/ankerl/exmars.htm) *(modified by myself)*

#### Front-end

* [skeleton](http://getskeleton.com/)
* [cash](http://kenwheeler.github.io/cash/)
* [font awesome](https://fortawesome.github.io/Font-Awesome)
* [chart.js](http://www.chartjs.org/)
* [sweat alert](http://t4t5.github.io/sweetalert/)
