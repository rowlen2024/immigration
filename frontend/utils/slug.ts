import { pinyin } from 'pinyin-pro';

export function generateSlugFromText(text: string): string {
  const words: string[] = [];
  let asciiWord = '';

  const flushAsciiWord = () => {
    if (!asciiWord) return;
    words.push(asciiWord.toLowerCase());
    asciiWord = '';
  };

  for (const char of text) {
    if (/^[A-Za-z0-9]$/.test(char)) {
      asciiWord += char;
      continue;
    }

    flushAsciiWord();

    const converted = pinyin(char, { toneType: 'none', type: 'array' })
      .join('')
      .toLowerCase();

    if (/^[a-z0-9]+$/.test(converted)) {
      words.push(converted);
    }
  }

  flushAsciiWord();

  return words
    .join('-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
}
