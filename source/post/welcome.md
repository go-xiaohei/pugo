```ini
title = Welcome
slug = welcome
desc = welcome to try pugo.static site generator
date = 2015-11-28
update_date = 2015-11-28
author = pugo
tags = pugo,welcome
```

When you read the blog, `pugo` is running successfully. Then enjoy your blog journey.

This blog is generated from file `source/welcome.md`. You can learn it and try to write your own blog article with following guide together.

#### blog meta

Blog's meta data, such as title, author, are created by first `ini` section with block **\`\`\`ini ..... \`\`\`**:

    ```ini
    title = Welcome to Pugo.Static
    slug = welcome-pugo-static
    date = 2015-11-08
    update_date = 2015-11-11
    author = pugo-robot
    author_email =
    author_url =
    tags = pugo,welcome
    ```

#### blog content

Content are from the second section as `markdown` format:

```markdown
When you read the blog, `pugo` is running successfully. Then enjoy your blog journey.

This blog is generated from file `source/welcome.md`. You can learn it and try to write your own blog article with following guide together.

...... (markdown content)
```

Just write content after blog meta, all words will be parsed as markdown content.
