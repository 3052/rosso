package amazon

const (
   DeviceID = "deviceID"
   // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
   HostATVExt    = "https://atv-ext.amazon.com"

   // the wrong DTID will fail the license request. if you change the DTID you
   // need to relog

   // no
   //DeviceTypeID = "A2SNKIF736WF4T" // com.amazon.amazonvideo.livingroom

   // yes
   DeviceTypeID = "A3NM0WFSU3DLT5" // sea_of_silence

   // maybe
   //DeviceTypeID = "A12GXV8XMS007S" // fire_tv_gen2
   //DeviceTypeID = "A1C66CX2XD756O" // fire_hd_8
   //DeviceTypeID = "A1KAXIG6VXSG8Y" // nvidia_shield: nvidia shield, unknown which one or if all
   //DeviceTypeID = "A1Q7QCGNMXAKYW" // fire_7_again: not sure the difference
   //DeviceTypeID = "A1RTAM01W29CUP" // pc_app
   //DeviceTypeID = "A1ZB65LA390I4K" // fire_hd_10
   //DeviceTypeID = "A265XOI9586NML" // fire_tv_stick_with_alexa
   //DeviceTypeID = "A2E0SNTXJVT7WK" // fire_tv: this is not the stick, this is the older stick-like diamond shaped one
   //DeviceTypeID = "A2GFL5ZMWNE0PX" // fire_tv_stick_4k: 4k fire tv stick
   //DeviceTypeID = "A2JKHJ0PX4J3L3" // fire_tv_cube: this is the STB-style big bulky cube
   //DeviceTypeID = "A2LWARUGJLBYEW" // fire_tv_stick_gen2
   //DeviceTypeID = "A2M4YX06LWP8WI" // fire_7
   //DeviceTypeID = "A32DOYMUN6DTXA" // echo_dot: echo dot Gen3
   //DeviceTypeID = "A38EHHIB10L47V" // fire_hd_8_again: not sure the difference
   //DeviceTypeID = "A3RBAYBE7VM004" // echo_studio: for audio stuff, this is probably the one to use
   //DeviceTypeID = "A43PXU4ZN2AL1"  // mobile_app
   //DeviceTypeID = "A71I8788P1ZV8"  // device_type
   //DeviceTypeID = "A7WXQPH584YP"   // echo: echo Gen2
   //DeviceTypeID = "ADVBD696BHNV5"  // fire_tv_stick_gen1: non-4k fire tv stick
   //DeviceTypeID = "AKPGW064GI9HE"  // fire_tv_stick_4k_gen3
   //DeviceTypeID = "AOAGZA014O5RE"  // browser: all browsers? all platforms?
   //DeviceTypeID = "AVU7CPPF2ZRAS"  // fire_hd_8_plus_2020

)
