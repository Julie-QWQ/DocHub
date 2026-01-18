import { computed, type Ref } from 'vue'
import { useForm, useField } from 'vee-validate'
import * as yup from 'yup'

/**
 * 登录表单验证
 */
export function useLoginForm(loginMethod?: Ref<'username' | 'email'>) {
  const schema = computed(() => {
    const usernameSchema = yup.string()
      .required('请输入用户名或邮箱')
      .min(3, '用户名至少3个字符')
      .max(50, '用户名最多50个字符')

    const passwordSchema = loginMethod?.value === 'email'
      ? yup.string().notRequired()
      : yup.string()
        .required('请输入密码')
        .min(6, '密码至少6个字符')
        .max(50, '密码最多50个字符')

    return yup.object({
      username: usernameSchema,
      password: passwordSchema
    })
  })

  const { handleSubmit, errors, isSubmitting } = useForm({
    validationSchema: schema
  })

  const { value: username } = useField('username')
  const { value: password } = useField('password')

  return {
    username,
    password,
    errors,
    isSubmitting,
    handleSubmit
  }
}

/**
 * 注册表单验证
 */
export function useRegisterForm() {
  const schema = yup.object({
    username: yup.string()
      .required('请输入用户名')
      .min(3, '用户名至少3个字符')
      .max(50, '用户名最多50个字符')
      .matches(/^[a-zA-Z0-9_]+$/, '用户名只能包含字母、数字和下划线'),
    email: yup.string()
      .required('请输入邮箱')
      .email('请输入有效的邮箱地址')
      .max(100, '邮箱最多100个字符'),
    password: yup.string()
      .required('请输入密码')
      .min(6, '密码至少6个字符')
      .max(50, '密码最多50个字符'),
    confirmPassword: yup.string()
      .required('请确认密码')
      .oneOf([yup.ref('password')], '两次输入的密码不一致'),
    real_name: yup.string()
      .required('请输入真实姓名')
      .min(2, '姓名至少2个字符')
      .max(50, '姓名最多50个字符'),
    major: yup.string()
      .required('请输入专业')
      .min(2, '专业至少2个字符')
      .max(100, '专业最多100个字符'),
    class: yup.string()
      .required('请输入班级')
      .min(2, '班级至少2个字符')
      .max(50, '班级最多50个字符')
  })

  const { handleSubmit, errors, isSubmitting } = useForm({
    validationSchema: schema
  })

  const { value: username } = useField('username')
  const { value: email } = useField('email')
  const { value: password } = useField('password')
  const { value: confirmPassword } = useField('confirmPassword')
  const { value: real_name } = useField('real_name')
  const { value: major } = useField('major')
  // 使用不同的方式处理 class 字段
  const classField = useField('class')

  return {
    username,
    email,
    password,
    confirmPassword,
    real_name,
    major,
    className: classField.value,
    errors,
    isSubmitting,
    handleSubmit
  }
}

/**
 * 修改密码表单验证
 */
export function useChangePasswordForm() {
  const schema = yup.object({
    old_password: yup.string()
      .required('请输入当前密码')
      .min(6, '密码至少6个字符')
      .max(50, '密码最多50个字符'),
    new_password: yup.string()
      .required('请输入新密码')
      .min(6, '密码至少6个字符')
      .max(50, '密码最多50个字符'),
    confirm_password: yup.string()
      .required('请确认新密码')
      .oneOf([yup.ref('new_password')], '两次输入的密码不一致')
  })

  const { handleSubmit, errors, isSubmitting } = useForm({
    validationSchema: schema
  })

  const { value: old_password } = useField('old_password')
  const { value: new_password } = useField('new_password')
  const { value: confirm_password } = useField('confirm_password')

  return {
    old_password,
    new_password,
    confirm_password,
    errors,
    isSubmitting,
    handleSubmit
  }
}
