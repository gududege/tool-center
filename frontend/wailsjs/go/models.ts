export namespace wails {
	
	export class CancelTaskResponse {
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CancelTaskResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	    }
	}
	export class DeleteTaskResponse {
	    success: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DeleteTaskResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	    }
	}
	export class DialogResponse {
	    path?: string;
	
	    static createFrom(source: any = {}) {
	        return new DialogResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class ParameterMappingDto {
	    field: string;
	    kind: string;
	    flag?: string;
	    style?: string;
	    separator?: string;
	    trueFlag?: string;
	    falseFlag?: string;
	    defaultValue?: any;
	
	    static createFrom(source: any = {}) {
	        return new ParameterMappingDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.field = source["field"];
	        this.kind = source["kind"];
	        this.flag = source["flag"];
	        this.style = source["style"];
	        this.separator = source["separator"];
	        this.trueFlag = source["trueFlag"];
	        this.falseFlag = source["falseFlag"];
	        this.defaultValue = source["defaultValue"];
	    }
	}
	export class ExecutionDefinitionDto {
	    exe: string;
	    workingDirectory?: string;
	    parameters?: ParameterMappingDto[];
	
	    static createFrom(source: any = {}) {
	        return new ExecutionDefinitionDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exe = source["exe"];
	        this.workingDirectory = source["workingDirectory"];
	        this.parameters = this.convertValues(source["parameters"], ParameterMappingDto);
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
	export class FileFilterDto {
	    displayName: string;
	    patterns: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileFilterDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.patterns = source["patterns"];
	    }
	}
	export class FormDefinitionDto {
	    schema?: number[];
	    uiSchema?: number[];
	
	    static createFrom(source: any = {}) {
	        return new FormDefinitionDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.uiSchema = source["uiSchema"];
	    }
	}
	export class NavigationDto {
	    group: string[];
	    group_cn?: string[];
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new NavigationDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.group = source["group"];
	        this.group_cn = source["group_cn"];
	        this.order = source["order"];
	    }
	}
	export class OutputEventDto {
	    taskId: string;
	    timestamp: string;
	    level: string;
	    source: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new OutputEventDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.taskId = source["taskId"];
	        this.timestamp = source["timestamp"];
	        this.level = source["level"];
	        this.source = source["source"];
	        this.message = source["message"];
	    }
	}
	
	export class PluginMetadataDto {
	    id: string;
	    name: string;
	    name_cn?: string;
	    description?: string;
	    description_cn?: string;
	    version?: string;
	    author?: string;
	    icon?: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginMetadataDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_cn = source["name_cn"];
	        this.description = source["description"];
	        this.description_cn = source["description_cn"];
	        this.version = source["version"];
	        this.author = source["author"];
	        this.icon = source["icon"];
	    }
	}
	export class PluginDto {
	    metadata: PluginMetadataDto;
	    navigation: NavigationDto;
	    form: FormDefinitionDto;
	    execution: ExecutionDefinitionDto;
	
	    static createFrom(source: any = {}) {
	        return new PluginDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.metadata = this.convertValues(source["metadata"], PluginMetadataDto);
	        this.navigation = this.convertValues(source["navigation"], NavigationDto);
	        this.form = this.convertValues(source["form"], FormDefinitionDto);
	        this.execution = this.convertValues(source["execution"], ExecutionDefinitionDto);
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
	export class PluginHistoryEntryDto {
	    timestamp: string;
	    label: string;
	    formData: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new PluginHistoryEntryDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.label = source["label"];
	        this.formData = source["formData"];
	    }
	}
	
	export class PluginSummaryDto {
	    id: string;
	    name: string;
	    name_cn?: string;
	    description?: string;
	    description_cn?: string;
	    version?: string;
	    icon?: string;
	    navigation: NavigationDto;
	
	    static createFrom(source: any = {}) {
	        return new PluginSummaryDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_cn = source["name_cn"];
	        this.description = source["description"];
	        this.description_cn = source["description_cn"];
	        this.version = source["version"];
	        this.icon = source["icon"];
	        this.navigation = this.convertValues(source["navigation"], NavigationDto);
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
	export class ReloadPluginsResponse {
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ReloadPluginsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	    }
	}
	export class RunPluginRequest {
	    pluginId: string;
	    formData: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new RunPluginRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pluginId = source["pluginId"];
	        this.formData = source["formData"];
	    }
	}
	export class RunPluginResponse {
	    taskId: string;
	
	    static createFrom(source: any = {}) {
	        return new RunPluginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.taskId = source["taskId"];
	    }
	}
	export class SettingsDto {
	    theme: string;
	    formTheme?: string;
	    pluginDirectory: string;
	    language?: string;
	    sidebarCollapsed: boolean;
	    sidebarSize: number;
	    bottomPanelSize: number;
	    bottomTab: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.formTheme = source["formTheme"];
	        this.pluginDirectory = source["pluginDirectory"];
	        this.language = source["language"];
	        this.sidebarCollapsed = source["sidebarCollapsed"];
	        this.sidebarSize = source["sidebarSize"];
	        this.bottomPanelSize = source["bottomPanelSize"];
	        this.bottomTab = source["bottomTab"];
	    }
	}
	export class SystemInfoDto {
	    appVersion: string;
	    buildTime: string;
	    goVersion: string;
	    os: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfoDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appVersion = source["appVersion"];
	        this.buildTime = source["buildTime"];
	        this.goVersion = source["goVersion"];
	        this.os = source["os"];
	    }
	}
	export class TaskDto {
	    id: string;
	    pluginId: string;
	    status: string;
	    createdAt: string;
	    startedAt?: string;
	    endedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pluginId = source["pluginId"];
	        this.status = source["status"];
	        this.createdAt = source["createdAt"];
	        this.startedAt = source["startedAt"];
	        this.endedAt = source["endedAt"];
	    }
	}
	export class ValidationErrorDto {
	    path: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidationErrorDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.message = source["message"];
	    }
	}
	export class ValidationResultDto {
	    valid: boolean;
	    errors?: ValidationErrorDto[];
	
	    static createFrom(source: any = {}) {
	        return new ValidationResultDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.errors = this.convertValues(source["errors"], ValidationErrorDto);
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

