# Reel-time programming project - Group ID 18
We have chosen to go with a peer-to-peer logic for our elevators. <br/>

## How start the elevators:

To start the elevator write ```go run main.go *ID*``` in the terminal window. <br/>
The *ID* you choose should be unique for every elevator you want to run together. <br/>

## General description of code:
The program contains four folders and a main file. These four folders are: <br/>
- driver_go_master: a driver containting logic connecting the code to the physical elevators functionality. <br/>
- elevator: initializes one elevator, and contains the logic behind the finite-state machine directing the elevator. <br/>
- hallassign: calculates which elevator should take an order, and assignes it to the most sufficient elevator. <br/>
- network: transmitts the elevator running on the computer, and receives the other elevators including it self. <br/>
