foo = "ENC[AES256_GCM,data:VWTe,iv:1jXojBlsfxaGel2OlCAY9i0CDaaVRCudGjuJlV90EA8=,tag:yzj4lEr4OaccODNka2shBg==,type:str]"

baz "qux" {
  quux = [
    "ENC[AES256_GCM,data:fw==,iv:Qykwql3Kahp2uoSgTLUv60HGniPU839v2i2zoDPuxJ4=,tag:ZYU3fHRnO+spNg0+qJAH/A==,type:int]",
    "ENC[AES256_GCM,data:MA==,iv:CDH2zKXaYMhJr4c/528Z/hFM2IfLxlvQSc1cgVH+WT4=,tag:HDpCS1qqOJDU64PCpUcVZQ==,type:int]",
  ]
  ref = "ENC[AES256_GCM,data:AVv1laLwAWU=,iv:E37lz0Sw+aFKC6sS5xtMELeu4a3DliYQR5z4Qs2Dox0=,tag:x/PC6Lfj7OpNdsployZLhg==,type:str]"
}

sops = {
  lastmodified = "2026-10-16T20:25:39Z"
  mac          = "ENC[AES256_GCM,data:tg9KJb7A0nHIb+SmlZiry1DDQop8JIiJtN405BLNs1oBD1y3zb1XJjve5k837K6Hm7xTVrb4gw/MgW5c32IyZdjTqEQv+FusdYZ3ElK3OcIFFUoNvbakOBnGBTNA+h28ViH5c6LaH8awVftfdguOxouODqMhvlD5a0mHLy38oAc=,iv:lct8P+0EXKLNg6dh6elZk4kM09dzaljC1fbnCc3vv24=,tag:E8Wu1fqMUSPiuQG4Yg1a4Q==,type:str]"
  pgp = [
    {
      created_at = "2026-10-16T20:25:39Z"
      enc        = "-----BEGIN PGP MESSAGE-----\n\nhQEMAyUpShfNkFB/AQf+MBQLTQghvVfi1n2Bu8F+0OP8WW7vx0ZKSUY19piKrh+G\nq8kfvpGkA+xHiQ3DGR034jeTwAcOBPnncJLF+omeeIUEO7frT9I4e7EMMBAKfubF\n5V/me5Oq85OhjvokzlxguTJeZvzgLll7KL4/vletIpUc2SHK9beK5dNiuVVSLjNc\nY4w9zM1t7SGQL/ueRmxCQBywAJsQbA8sclQHYI/0+C/fmjm1qEcYKdLwYDVCXSKB\nDOFhL4acuKu+eMljpxMe4vJ+dp+CQqF59WOKjYicDDSplYc2LVw9i7addD6NRRJJ\nQY7mfBFbATBfDyAAQdRINrbWpS05o/DdE41IQ/hl/tJeAbTDlkI0MI4U8p+LrERa\n6+VyeDX7tC/MNKzSznutBAbYnVRz4sZ4G5EBPJBVpSCcPEa8hmB+zcMaehjGwnFL\ndPXDtRMR4nUgzHi9Iha4OJFo5D+N5F2pFWFyqw1hyg==\n=0/eM\n-----END PGP MESSAGE-----"
      fp         = "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"
    },
  ]
  unencrypted_suffix = "_unencrypted"
  version            = "3.13.3"
}
