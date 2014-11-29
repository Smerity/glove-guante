import sys

import numpy as np

if __name__ == '__main__':
  print 'Loading word vectors...'
  wordvecs = None
  wordlist = []
  for i, line in enumerate(sys.stdin):
    word, vec = line.strip().split(' ', 1)
    vec = map(float, vec.split())
    if wordvecs is None:
      wordvecs = np.ones((400000, len(vec)), dtype=np.float)
    wordvecs[i] = vec
    wordlist.append(word)
  words = dict((k, wordvecs[v]) for v, k in enumerate(wordlist))

  tests = [('he', words['he']), ('she', words['she'])]
  tests = [
      ('athens-greece+berlin', words['athens'] - words['greece'] + words['berlin']),
      ('sydney-australia+berlin', words['sydney'] - words['australia'] + words['berlin']),
      ('australia-sydney+germany', words['australia'] - words['sydney'] + words['berlin']),
      ('king-male+female', words['king'] - words['male'] + words['female']),
      ('king-man+woman', words['king'] - words['man'] + words['woman']),
      ('queen-female+male', words['queen'] - words['female'] + words['male']),
      ('queen-woman+man', words['queen'] - words['woman'] + words['man']),
      ('plane-air+rail', words['train'] - words['air'] + words['rail']),
  ]
  for test, tvec in tests:
    results = []
    print '=-=-' * 10
    print 'Testing {}'.format(test)
    res = np.dot(wordvecs, tvec) / np.linalg.norm(tvec) / np.linalg.norm(wordvecs, axis=1)
    results = zip(res, wordlist)
    print '\n'.join([w for _, w in sorted(results, reverse=True)[:20]])
