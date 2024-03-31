package common

var (
	Version   = "v1.0.0"
	GitHash   = ""
	BuildTime = ""
	GoVersion = ""
	Banner    = `

                                                                      
@@@  @@@  @@@@@@@   @@@@@@@              @@@@@@@   @@@@@@   @@@@@@@   
@@@  @@@  @@@@@@@@  @@@@@@@@             @@@@@@@  @@@@@@@@  @@@@@@@@  
@@!  @@@  @@!  @@@  @@!  @@@               @@!    @@!  @@@  @@!  @@@  
!@!  @!@  !@!  @!@  !@!  @!@               !@!    !@!  @!@  !@!  @!@  
@!@  !@!  @!@  !@!  @!@@!@!   @!@!@!@!@    @!!    @!@!@!@!  @!@@!@!   
!@!  !!!  !@!  !!!  !!@!!!    !!!@!@!!!    !!!    !!!@!!!!  !!@!!!    
!!:  !!!  !!:  !!!  !!:                    !!:    !!:  !!!  !!:       
:!:  !:!  :!:  !:!  :!:                    :!:    :!:  !:!  :!:       
::::: ::   :::: ::   ::                     ::    ::   :::   ::       
 : :  :   :: :  :    :                      :      :   : :   :                     
                                              
A simple tap tunnel encrypted over UDP, similar to VXLAN.
`
)

func PrintVersion() {
	println(Banner)
	println("Version: ", Version)
	println("GitHash: ", GitHash)
	println("BuildTime: ", BuildTime)
	println("GoVersion: ", GoVersion)
}
