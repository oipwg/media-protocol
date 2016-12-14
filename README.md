# Media-Protocol

[![Build Status][TravisSVG]][TravisLink] [![Coverage Status][CoverallsSVG]][CoverallsLink] [![Go Report Card][GoReportCardSVG]][GoReportCardLink]

## OIP-041 JSON Standards
The following are the current OIP-041 JSON Standards.

### Publish Artifact
```javascript
{  
    "oip-041":{  
        "artifact":{  
            "publisher":"$PublisherAddress",
            "timestamp":1234567890,
            "type":"$ArtifactType",
            "info":{  
                "title":"$ArtifactTitle",
                "description":"$ArtifactDescription",
                "year":1234,
                "extraInfo":{  
                    "artist":"$Creator",
                    "company":"$Distributor",
                    "composers":[  
                        "$Composer1",
                        "$Composer2"
                    ],
                    "copyright":"",
                    "usageProhibitions":"",
                    "usageRights":"",
                    "tags":[  
                        "$ArtifactTag1",
                        "$ArtifactTag2"
                    ]
                }
            },
            "storage":{  
                "network":"IPFS",
                "location":"$IPFSAddress",
                "files":[  
                    {  
                        "disallowBuy":true,
                        "dname":"$DisplayName",
                        "duration":123,
                        "fname":"$FileName",
                        "fsize":123,
                        "minPlay":"$minPlayPriceUSD",
                        "sugPlay":"$suggestedPlayPriceUSD",
                        "promo":"$CutForPromoterSales",
                        "retail":"$CutForPlatformSales",
                        "ptpFT":12,
                        "ptpDT":34,
                        "ptpDA":56,
                        "type":"$MediaType",
                        "tokenlyID":"$SongTokenlyID"
                    },
                    {  
                        "dissallowPlay":true,
                        "dname":"$DisplayName",
                        "duration":123,
                        "fname":"$FileName",
                        "fsize":123,
                        "minBuy":"$minBuyPriceUSD",
                        "sugBuy":"$suggestedBuyPriceUSD",
                        "promo":"$CutForPromoterSales",
                        "retail":"$CutForPlatformSales",
                        "type":"$MediaType",
                        "tokenlyID":"$SongTokenlyID"
                    },
                    {  
                        "dname":"$DisplayName",
                        "duration":123,
                        "fname":"$FileName",
                        "fsize":123,
                        "minPlay":"$minPlayPriceFiat",
                        "sugPlay":"$suggestedPlayPriceFiat",
                        "minBuy":"$minBuyPriceFiat",
                        "sugBuy":"$suggestedBuyPriceFiat",
                        "promo":"$CutForPromoterSales",
                        "retail":"$CutForPlatformSales",
                        "ptpFT":12,
                        "ptpDT":34,
                        "ptpDA":56,
                        "type":"$MediaType",
                        "tokenlyID":"$SongTokenlyID"
                    },
                    {  
                        "dname":"Cover Art",
                        "fname":"$CoverArtFilename",
                        "fsize":123,
                        "type":"coverArt",
                        "storage":{  
                            "network":"HTTP",
                            "location":"$ThumbnailURL"
                        }
                    }
                ]
            },
            "payment":{  
                "fiat":"$fiat_id",
                "scale":"1000:1",
                "sugTip":[  
                    123,
                    123,
                    123
                ],
                "tokens":{  
                    "btc":"$BitcoinAddress",
                    "early":"",
                    "mtmcollector":"",
                    "mtmproducer":"",
                    "happybirthdayep":"",
                    "ltbcoin":""
                }
            },
            "varID":""
        },
        "signature":"$IPFSAddress-$PublisherAddress-$timestamp"
    }
}
```
### Edit Artifact
```javascript
{  
    "oip-041":{  
        "editArtifact":{  
            "txid":"a449b0f6a601e503e7b4fdc0ada47f55a8b2f98feb2fdb044f7a92d971ff0456",
            "timestamp":1481420001,
            "patch":{  
                "add":[  
                    {  
                        "path":"/payment/tokens/mtcproducer",
                        "value":""
                    }
                ],
                "replace":[  
                    {  
                        "path":"/storage/files/3/fname",
                        "value":"birthdayepFirst.jpg"
                    },
                    {  
                        "path":"/storage/files/3/dname",
                        "value":"Cover Art 2"
                    },
                    {  
                        "path":"/info/title",
                        "value":"Happy Birthday"
                    },
                    {  
                        "path":"/timestamp",
                        "value":1481420001
                    }
                ],
                "remove":[  
                    {  
                        "path":"/payment/tokens/mtmproducer"
                    },
                    {  
                        "path":"/storage/files/0/sugBuy"
                    }
                ]
            }
        }
    }
}
```
### Transfer Artifact
```javascript
{  
    "oip-041":{  
        "transferArtifact":{  
            "txid":"$artifactID",
            "to":"$newPublisherAddress",
            "from":"$oldPublisherAddress",
            "timestamp":1234567890
        },
        "signature":"sign($artifactID-$newPublisherAddress-$oldPublisherAddress-$timestamp)"
    }
}
```
### Deactivate Artifact
```javascript
{  
    "oip-041":{  
        "deactivateArtifact":{  
            "txid":"96bad8e17f908da4c695c58b0f843a03928e338b361b3035eda16a864eafc3a2",
            "timestamp":1481697196
        },
        "signature":"H8wPKDHSTrrJIY4RVzoWuWlt5Fta4PaWbWNQvrGt9hRyBFrB2YDoVnhctC4V08KnGMD/CZKbA4cPKysStVq8jyE="
    }
}
```

### Multipart Data
You should not have to use the Multipart Data format, however should you need to reference it, here it is. When the JSON to be published is larger than 528 characters, the JSON gets split up into multiple parts. Each of these parts are then submitted as Transaction comments to the Florincoin Blockchain. Here is an example of an OIP 6 part artifact.

This data is formatted as follows:
```javascript
alexandria-media-multipart($partNumber, $numberOfParts, $publisherAddress, $firstPartTXID, $signature):$choppedStringData
```
```javascript
alexandria-media-multipart(0,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,0000000000000000000000000000000000000000000000000000000000000000,IBvdL1xJhvk2NIs7ckwsmK4hGGI2rnhgYwbTa6zy/FF1TxFyLuiv2fKTZYf7nmK0bHX0prUv4pl/CU/ZErvleW4=):{"oip-041":{"artifact":{"publisher":"FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q","timestamp":1481420000,"type":"music","info":{"title":"Happy Birthday EP","description":"this is the second organically grown, gluten free album released by Adam B. Levine - contact adam@tokenly.co
```
```javascript
alexandria-media-multipart(1,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,IBnj6xxykNf3ZDpidg8dk4ioERFU3Gj2tKQ3dxFAXeIQB3gPibrWF5b4g4PIR8KimwqqmqDQ77PF4dApAhuXze4=):m with questions or comments or discuss collaborations.","year":"2016","extraInfo":{"artist":"Adam B. Levine","company":"","composers":["Adam B. Levine"],"copyright":"","usageProhibitions":"","usageRights":"","tags":[]}},"storage":{"network":"IPFS","location":"QmPukCZKe
```
```javascript
alexandria-media-multipart(2,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,IMJbM7xMAVl/XWN0KpJiid/LADx+HdDNJUdiUkxm4JFyRKGCkF3VBt6cTUJ50YT3HO0heNMBCyh3HGFiunQWqis=):JD4KZFtstpvrguLaq94rsWfBxLU1QoZxvgRxA","files":[{"dname":"Skipping Stones","fame":"1 - Skipping Stones.mp3","fsize":6515667,"type":"album track","duration":1533.603293,"sugPlay":100,"minPlay":null,"sugBuy":750,"minBuy":500,"promo":10,"retail":15,"ptpFT":10,"ptpDT":20,"p
```
```javascript
alexandria-media-multipart(3,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,IBdKG047JgK2XyWX86AFf9n1yT+QTtPJjgOnofP74wwBHZBYT4gJQeSNOfToIrerNkcXGr9zX+N1nTVhKLuscx0=):tpDA":50},{"dname":"Lessons","fame":"2 - Lessons with intro.mp3","fsize":6515667,"type":"album track","duration":1231.155243,"disallowPlay":1,"sugBuy":750,"minBuy":500,"promo":10,"retail":15,"ptpFT":10,"ptpDT":20,"ptpDA":50},{"dname":"Born to Roam","fame":"3 - Born to R
```
```javascript
alexandria-media-multipart(4,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,H7KGidUTMG+6xBwudl1EXeFcBvGy9+UHnd9vrWC7ETTBCioioE0pphkozUgbZQx+jIldlEPMBnVuqox383P8nLI=):oam.mp3","fsize":6515667,"type":"album track","duration":2374.550714,"sugPlay":100,"minPlay":50,"disallowBuy":1,"promo":10,"retail":15,"ptpFT":10,"ptpDT":20,"ptpDA":50},{"dname":"Cover Art","fname":"birthdayepFINAL.jpg","type":"coverArt","disallowBuy":1}]},"payment":{"f
```
```javascript
alexandria-media-multipart(5,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,INaPk7aMwksd9rzRXtCqND9RJYZPbUagosessK+b4D+JUfga/gT1gU25lvTs2hLZkcoqfVXGqlRsOrZ2agGYw3M=):iat":"USD","scale":"1000:1","sugTip":[5,50,100],"tokens":{"mtmcollector":"","mtmproducer":"","happybirthdayep":"","early":"","ltbcoin":"","btc":"1GMMg2J5iUKnDf5PbRr9TcKV3R6KfUiB55"}}},"signature":"H3XC/u9qz9pUP5g1+dyWUSR2euKFH3DWKd8hTdFINURvZvcdE7Q2hnNJa9QOuunCD1CPiVMOV
```
```javascript
alexandria-media-multipart(6,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,5d0eb0bfb05815567717ec1d5b72c92c8bcf8d30c48785d6449970bb32a9c07b,HwRHpvyi99EM0xtA68FGLWJpd4sls/z6zNAjQh65OnQhRp19mSZNqoheYdw6a4QReUd0I0iBvMt0udgrIXLuE6Y=):q+8m+NcgMQTw60="}}
```
## License

This project uses the [MIT] License.



[TravisSVG]: https://travis-ci.org/dloa/media-protocol.svg?branch=master
[TravisLink]: https://travis-ci.org/dloa/media-protocol

[CoverallsSVG]: https://coveralls.io/repos/github/dloa/media-protocol/badge.svg?branch=master
[CoverallsLink]: https://coveralls.io/github/dloa/media-protocol?branch=master

[GoReportCardSVG]: https://goreportcard.com/badge/github.com/dloa/media-protocol
[GoReportCardLink]: https://goreportcard.com/report/github.com/dloa/media-protocol

[MIT]:LICENSE.md