import { useI18n } from 'vue-i18n'

export function useLocaleText() {
  const { locale } = useI18n()

  function localText(textObj: { zh: string; en: string } | undefined | null): string {
    if (!textObj) return ''
    return locale.value === 'zh-CN' ? textObj.zh : textObj.en
  }

  return { localText }
}
