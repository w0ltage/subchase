package version

import "fmt"

func ShowBanner(version string) {
    fmt.Printf(
`               __         __                  
   _______  __/ /_  _____/ /_  ____ _________ 
  / ___/ / / / __ \/ ___/ __ \/ __ %c/ ___/ _ \
 (__  ) /_/ / /_/ / /__/ / / / /_/ (__  )  __/
/____/\__,_/_.___/\___/_/ /_/\__,_/____/\___/  %v

`, '`', version)
}

