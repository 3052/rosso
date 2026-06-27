package amazon

const (
   // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
   HostATVExt    = "https://atv-ext.amazon.com"

   // the wrong DTID will fail the license request. if you change the DTID you
   // need to relog

   // time.is/Unix_time
   DeviceName = "1782583253"

   //DeviceID = "deviceID"
   DeviceID = "uuidd9b2d63da09b46c1952a8c7b3e1f0a4d"

   // state 0 untested

   // state 1 pass codepair
   //DeviceTypeID = "A12GXV8XMS007S" // fire_tv_gen2
   //DeviceTypeID = "A1C66CX2XD756O" // fire_hd_8
   //DeviceTypeID = "A1Q7QCGNMXAKYW" // fire_7_again: not sure the difference
   //DeviceTypeID = "A1ZB65LA390I4K" // fire_hd_10
   //DeviceTypeID = "A265XOI9586NML" // fire_tv_stick_with_alexa
   //DeviceTypeID = "A2E0SNTXJVT7WK" // fire_tv: this is not the stick, this is the older stick-like diamond shaped one
   //DeviceTypeID = "A2GFL5ZMWNE0PX" // fire_tv_stick_4k: 4k fire tv stick
   //DeviceTypeID = "A2JKHJ0PX4J3L3" // fire_tv_cube: this is the STB-style big bulky cube
   //DeviceTypeID = "A2LWARUGJLBYEW" // fire_tv_stick_gen2
   //DeviceTypeID = "A2M4YX06LWP8WI" // fire_7
   //DeviceTypeID = "A38EHHIB10L47V" // fire_hd_8_again: not sure the difference
   //DeviceTypeID = "A7WXQPH584YP"   // echo: echo Gen2
   //DeviceTypeID = "ADVBD696BHNV5"  // fire_tv_stick_gen1: non-4k fire tv stick
   //DeviceTypeID = "AKPGW064GI9HE"  // fire_tv_stick_4k_gen3
   //DeviceTypeID = "AVU7CPPF2ZRAS"  // fire_hd_8_plus_2020

   // state 2 pass register
   //DeviceTypeID = "A1RTAM01W29CUP" // pc_app
   //DeviceTypeID = "A3EFHJ9BGBJ8L2" // LegacyRefreshToken
   //DeviceTypeID = "A43PXU4ZN2AL1"  // mobile_app

   // state 3 pass license FHD
   //DeviceTypeID = "A1KAXIG6VXSG8Y" // nvidia_shield: nvidia shield, unknown which one or if all
   //DeviceTypeID = "A2HYAJ0FEWP6N3" // MTC soc MStar_T22 Android9/10/11
   //DeviceTypeID = "A2RGJ95OVLR12U" // Hisense soc MTK9602 Vidaa4.0+ Linux
   //DeviceTypeID = "A2SNKIF736WF4T" // com.amazon.amazonvideo.livingroom
   //DeviceTypeID = "A394LFCMDJ1B8R" // LG UR Series soc LM21ANN_W23 webOS23/24
   //DeviceTypeID = "A3REWRVYBYPKUM" // Hisense HE55A7000EUWTS soc MSD6886 Vidaa4.0 Linux
   //DeviceTypeID = "A71I8788P1ZV8" // LG Mediatek soc cert_model_W3.0 webOS3.0
   //DeviceTypeID = "AOAGZA014O5RE" // Chromium based browsers

   // state 4 pass license UHD
   DeviceTypeID = "A3NM0WFSU3DLT5" // sea_of_silence
)

/*
98
A2M4YX06LWP8WI

88
A2E0SNTXJVT7WK

87
A12GXV8XMS007S

84
A1C66CX2XD756O

77
A2GFL5ZMWNE0PX

75
A2JKHJ0PX4J3L3
A2LWARUGJLBYEW

67
A38EHHIB10L47V

65
A265XOI9586NML

58
A1Q7QCGNMXAKYW

43
A1ZB65LA390I4K
*/
