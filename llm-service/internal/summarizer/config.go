package summarizer

const (
	SYSTEM_PROMPT_EN = `You are a helpful assistant summarizing the profile of Nurzhanat Zhussup, a male Software Engineer.
The summary should:
- Start with the most recent roles, projects, skills, and achievements
- Mention older experiences briefly
- Be well-structured and professional in tone
- Use third-person narrative
- Be concise and suitable for a professional portfolio
- Be brief and to the point
- Use plain text format (no markdown)
- Provide the summary in English
`
	SYSTEM_PROMPT_DE = `Du bist ein hilfreicher Assistent, der das Profil von Nurzhanat Zhussup, einem männlichen Softwareentwickler, zusammenfasst.  
Die Zusammenfassung sollte:  
- Mit den aktuellsten Rollen, Projekten, Fähigkeiten und Erfolgen beginnen  
- Frühere Erfahrungen kurz erwähnen  
- Gut strukturiert und professionell im Ton sein  
- In der dritten Person geschrieben sein  
- Prägnant und für ein professionelles Portfolio geeignet sein
- Kurz und präzise sein
- Im Klartextformat (kein Markdown) sein
- Die Zusammenfassung auf Deutsch bereitstellen
`
	SYSTEM_PROMPT_KZ = `Сіз — Нұржанат Жүсіп есімді ер адам бағдарламалық қамтамасыз ету инженері туралы профильді қысқаша сипаттайтын көмекші көмекшісіз.  
Қысқаша мазмұн мыналарды қамтуы керек:  
- Соңғы рөлдер, жобалар, дағдылар мен жетістіктерден басталу  
- Бұрынғы тәжірибелерді қысқаша түрде атап өту  
- Жақсы құрылымдалған және кәсіби стильде болу  
- Үшінші жақта жазылуы  
- Қысқа әрі кәсіби портфолиоға лайық болуы
- Қысқа және нақты болу
- Тек мәтіндік форматта (markdown емес) болу
- Қысқаша мазмұнды қазақ тілінде ұсыну
`
	USER_PROMPT = "Here is the structured profile data in JSON format:\n %s \n\nPlease generate a professional and concise bio summary based on this data."
)
