# CoAP_Server
This Project is a fast Proof of Concept to test the abilities of CoAP for a
GPS Tracker. All locations are stored in a global array and don't persist
For this it is implementing a single path for the CoAP server, `/locations`
over this route the server accepts a JSON containing of an imei, latitude 
and longitude values. For the as for the timestamp value, it gets added on
server side to sort the values. 

``` json
{
    "imei" : "ASDFGHJKL",
    "latitude" : 1,2345,
    "longitude" : 2,3456,
    "timestamp" : "0001-01-01T00:00Z"
}
```

For Debugging this project also integrates a simple REST API to get the device 
locations and add device locations. Also it contains a **VERY** simple htmx interface
to show the devicelocations. 

## Testing the CoAP path
Since there is no way to use the Copper extension for Firefox or Chrome, the tool of 
choise is this [coap-cli](https://github.com/mainflux/coap-cli) tool.

Example: `./coap-cli post /location -d "{\"imei\": \"aaaaaaaab\", \"latitude\": 1.234, \"longitude\": 2.345}"`
`
