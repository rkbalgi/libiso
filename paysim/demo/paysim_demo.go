package demo

var Demo_Specs string =`

{ "Specs":[{"SpecName": "ISO8583_1 v1 (ASCII)", "Fields": [{"BitPosition":0,"Name": "Message Type", "Type":"Fixed", "Attrs":"false;4;ascii;default"},
                                                           {"BitPosition":0,"Name": "Bitmap", "Type":"Bitmapped", "Attrs":"false;_;_;default","Children":
                                                          [{"BitPosition":2,"Name": "PAN", "Type":"Variable", "Attrs":"false;2;ascii;ascii;default"},
                                                           {"BitPosition":3,"Name": "Processing Code", "Type":"Fixed", "Attrs":"false;6;ascii;default"}, 
                                                           {"BitPosition":4,"Name": "Txn Amount", "Type":"Fixed", "Attrs":"false;12;ascii;default"}, 
                                                           {"BitPosition":11,"Name": "STAN", "Type":"Fixed", "Attrs":"false;6;ascii;default"}, 
                                                           {"BitPosition":14,"Name": "Expiry Date", "Type":"Fixed", "Attrs":"false;4;ascii;default"}, 
                                                           {"BitPosition":35,"Name": "Track 2", "Type":"Variable", "Attrs":"false;2;ascii;ascii;default"},
                                                           {"BitPosition":38,"Name": "Approval Code", "Type":"Fixed", "Attrs":"false;6;ascii;default"},
                                                           {"BitPosition":39,"Name": "Action Code", "Type":"Fixed", "Attrs":"false;3;ascii;default"},
                                                           {"BitPosition":52,"Name": "PIN Data", "Type":"Fixed", "Attrs":"false;8;binary;_"},
                                                           {"BitPosition":55,"Name": "ICC Data", "Type":"Variable", "Attrs":"false;3;ascii;binary;default"},
                                                           {"BitPosition":64,"Name": "MAC(1)", "Type":"Fixed", "Attrs":"false;8;binary;_"},
                                                           {"BitPosition":96,"Name": "Key Management Data", "Type":"Variable", "Attrs":"false;2;ascii;binary;default"},
                                                           {"BitPosition":128,"Name": "MAC(2)", "Type":"Fixed", "Attrs":"false;8;binary;_"}
                                                          ]
                                                        }  
                                             ]
          },
          {"SpecName": "ISO8583_1 v1 (EBCDIC)", "Fields": [{"Name": "Message Type", "Type":"Fixed", "Attrs":"false;4;ebcdic;default"},
                                                           {"Name": "Bitmap", "Type":"Bitmapped", "Attrs":"false;_;_;default","Children":
                                                          [{"BitPosition": 2,"Name": "PAN", "Type":"Variable", "Attrs":"false;2;ebcdic;ebcdic;default"},
                                                           {"BitPosition": 3,"Name": "Processing Code", "Type":"Fixed", "Attrs":"false;6;ebcdic;default"},
                                                           {"BitPosition": 4,"Name": "Txn Amount", "Type":"Fixed", "Attrs":"false;12;ebcdic;default"},                                                           
                                                           {"BitPosition": 11,"Name": "STAN", "Type":"Fixed", "Attrs":"false;6;ebcdic;default"}, 
                                                           {"BitPosition": 14,"Name": "Expiry Date", "Type":"Fixed", "Attrs":"false;4;ebcdic;default"}, 
                                                           {"BitPosition": 35,"Name": "Track 2", "Type":"Variable", "Attrs":"false;2;ebcdic;ebcdic;default"},
                                                           {"BitPosition": 38,"Name": "Approval Code", "Type":"Fixed", "Attrs":"false;6;ebcdic;default"},
                                                           {"BitPosition": 39,"Name": "Action Code", "Type":"Fixed", "Attrs":"false;3;ebcdic;default"},
                                                           {"BitPosition": 52,"Name": "PIN Data", "Type":"Fixed", "Attrs":"false;8;binary;_"},
                                                           {"BitPosition": 55,"Name": "ICC Data", "Type":"Variable", "Attrs":"false;3;ebcdic;binary;default"},
                                                           {"BitPosition": 64,"Name": "MAC(1)", "Type":"Fixed", "Attrs":"false;8;binary;_"},
                                                           {"BitPosition": 96,"Name": "Key Management Data", "Type":"Variable", "Attrs":"false;2;ebcdic;binary;?"},
                                                           {"BitPosition": 128,"Name": "MAC(2)", "Type":"Fixed", "Attrs":"false;8;binary;_"}
                                                          ]
                                                        }  
                                             ]
          }
          ]
}          `;
