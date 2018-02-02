# Radio Hub (RHub)

Provide a hub supporting many different protocols frequently used to control Ham Radio Equipments

for example :
   - PTY/Serial port to emulate a RIG CAT interface
   - TCP server to provide serial over lan type of interface (similar to ser2net)
   - Serial Port interface with RIG
   - Generic transceiver model to multiplex between the different clients.
   - Radio specific protocol adapter (CAT, CIV, ...)
   
The architecture is loosely based on the SYSV stream concepts allowing attaching modules in different configurations like :

    +-----------------------------+
    |       Generic Radio         |
    +-----------------------------+
         |          |          |
    +---------+  +------+  +------+
    |CMDBuffer|  |FT991A|  |FT991A|
    +---------+  +------+  +------+
        |            |         |
    +------+ +-----------+ +------+
    |  PTY | | CMDBuffer | |Serial|
    +------+ +-----------+ +------+
                   |
             +-----------+
             | TCPServer |
             +-----------+
