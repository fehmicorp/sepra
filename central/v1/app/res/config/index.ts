export const URL = {
    "domain":"fehmicorp.in",
    "web": "https://fehmicorp.in",
    "dist" : "https://dist.fehmicorp.in/"
} as const;

export const S3Bucket = {
    "V1":{
        "Assets":{
            "logo64":"v1/assets/logo/64.png",
            "logo192":"v1/assets/logo/192.png",
            "logo512":"v1/assets/logo/512.png"
        }
    }
}

export const Config = {
    "Name": "SEPRA",
    "Logo":{
        "64": `${URL.dist}/${S3Bucket.V1.Assets.logo64}`,
        "192": `${URL.dist}/${S3Bucket.V1.Assets.logo192}`,
        "512": `${URL.dist}/${S3Bucket.V1.Assets.logo512}`,
    },    
    "Company": {
        "LegalName": "Fehmi Corporation",
    }
} as const;

export type UrlDetails = typeof URL;
export type AppDetails = typeof Config;