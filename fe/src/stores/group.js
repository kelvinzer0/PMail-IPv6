import { ref } from 'vue'
import { defineStore } from 'pinia'
import { lang } from '../i18n/i18n';

const useGroupStore = defineStore('group', () => {
  const tag = ref("")
  const name = ref(lang.inbox)
  const searchQuery = ref('')
  return { tag, name, searchQuery }
})

export default useGroupStore