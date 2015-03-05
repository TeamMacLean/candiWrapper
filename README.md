# CandiWrapper


<img align="right" height="300" src="https://raw.githubusercontent.com/wookoouk/candiWrapper/master/logo.png">

> An interface between Galaxy and CandiSNP

CandiWrapper is a tool that sits between Galaxy and CandiSNP to handle interaction between the two.
When Galaxy posts snpEff data to it, it save the data to a redis database and return a url to the user, when the user follows the link it opens CandiSNP which pulls the data to be rendered.




## License

Copyright (c) 2015 The Sainsbury Laboratory and other contributors

Licensed under the MIT License
