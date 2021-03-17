import { ParamsType, ErrorType, ParamFields } from '../types/form'

const validate = (values: ParamsType) => {
  const errors: ErrorType = {};

  const requiredFields = [
    ParamFields.email,
    ParamFields.fullName,
    ParamFields.company,
  ];

  requiredFields.forEach(field => {
    if (!values[field]) {
      errors[field] = true;
    }
  });

  if (
    values.email &&
    !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(values.email)
  ) {
    errors[ParamFields.email] = true;
  }

  return errors;
}

export default validate
