export namespace main {
	
	export class AppConfig {
	    AppName: string;
	    Description: string;
	    Icon: string;
	    Tagline: string;
	    Title: string;
	    Width: number;
	    Height: number;
	    Min: boolean;
	    Max: boolean;
	    Quit: boolean;
	    Version: string;
	    Domain: string;
	    installDir: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.Description = source["Description"];
	        this.Icon = source["Icon"];
	        this.Tagline = source["Tagline"];
	        this.Title = source["Title"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	        this.Min = source["Min"];
	        this.Max = source["Max"];
	        this.Quit = source["Quit"];
	        this.Version = source["Version"];
	        this.Domain = source["Domain"];
	        this.installDir = source["installDir"];
	    }
	}
	export class InstallConfig {
	    isPrivate: boolean;
	    fqdn?: string;
	    port?: string;
	    apiKey?: string;
	    licenseData?: string;
	    installPath: string;
	    installAsService: boolean;
	    enableAutostart: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InstallConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isPrivate = source["isPrivate"];
	        this.fqdn = source["fqdn"];
	        this.port = source["port"];
	        this.apiKey = source["apiKey"];
	        this.licenseData = source["licenseData"];
	        this.installPath = source["installPath"];
	        this.installAsService = source["installAsService"];
	        this.enableAutostart = source["enableAutostart"];
	    }
	}
	export class SysPayload {
	    isAdmin: boolean;
	    targetArch: string;
	
	    static createFrom(source: any = {}) {
	        return new SysPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isAdmin = source["isAdmin"];
	        this.targetArch = source["targetArch"];
	    }
	}
	export class SystemPayload {
	    hostname: string;
	    targetOs: string;
	    targetArch: string;
	    targetBuild: string;
	    agentVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.targetOs = source["targetOs"];
	        this.targetArch = source["targetArch"];
	        this.targetBuild = source["targetBuild"];
	        this.agentVersion = source["agentVersion"];
	    }
	}

}

