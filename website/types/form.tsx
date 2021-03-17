export enum ParamFields {
  email = 'email',
  fullName = 'fullName',
  company = 'company'
}

export interface ParamsType {
  [ParamFields.email]: string | null;
  [ParamFields.fullName]: string | null;
  [ParamFields.company]: string | null;
}

export interface ErrorType {
  [ParamFields.email]?: boolean;
  [ParamFields.fullName]?: boolean;
  [ParamFields.company]?: boolean;
}
