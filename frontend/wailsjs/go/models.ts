export namespace cpu {
	
	export class InfoStat {
	    cpu: number;
	    vendorId: string;
	    family: string;
	    model: string;
	    stepping: number;
	    physicalId: string;
	    coreId: string;
	    cores: number;
	    modelName: string;
	    mhz: number;
	    cacheSize: number;
	    flags: string[];
	    microcode: string;
	
	    static createFrom(source: any = {}) {
	        return new InfoStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.vendorId = source["vendorId"];
	        this.family = source["family"];
	        this.model = source["model"];
	        this.stepping = source["stepping"];
	        this.physicalId = source["physicalId"];
	        this.coreId = source["coreId"];
	        this.cores = source["cores"];
	        this.modelName = source["modelName"];
	        this.mhz = source["mhz"];
	        this.cacheSize = source["cacheSize"];
	        this.flags = source["flags"];
	        this.microcode = source["microcode"];
	    }
	}

}

export namespace helpers {
	
	export class CPUInformation {
	    CPUPercent: number[];
	    CPUCores: number;
	    CPUInfo: cpu.InfoStat[];
	    CPUFrequency: number[];
	    CPUTemperature: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new CPUInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CPUPercent = source["CPUPercent"];
	        this.CPUCores = source["CPUCores"];
	        this.CPUInfo = this.convertValues(source["CPUInfo"], cpu.InfoStat);
	        this.CPUFrequency = source["CPUFrequency"];
	        this.CPUTemperature = source["CPUTemperature"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HostInformation {
	    MajorInfo?: host.InfoStat;
	    Uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new HostInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MajorInfo = this.convertValues(source["MajorInfo"], host.InfoStat);
	        this.Uptime = source["Uptime"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace host {
	
	export class InfoStat {
	    hostname: string;
	    uptime: number;
	    bootTime: number;
	    procs: number;
	    os: string;
	    platform: string;
	    platformFamily: string;
	    platformVersion: string;
	    kernelVersion: string;
	    kernelArch: string;
	    virtualizationSystem: string;
	    virtualizationRole: string;
	    hostId: string;
	
	    static createFrom(source: any = {}) {
	        return new InfoStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.uptime = source["uptime"];
	        this.bootTime = source["bootTime"];
	        this.procs = source["procs"];
	        this.os = source["os"];
	        this.platform = source["platform"];
	        this.platformFamily = source["platformFamily"];
	        this.platformVersion = source["platformVersion"];
	        this.kernelVersion = source["kernelVersion"];
	        this.kernelArch = source["kernelArch"];
	        this.virtualizationSystem = source["virtualizationSystem"];
	        this.virtualizationRole = source["virtualizationRole"];
	        this.hostId = source["hostId"];
	    }
	}

}

