Interesting that none of the tools besides Zarf even attempt to do concurrent writing. Tools will concurrently write layers but will not create a go routine to do the image pulling in parallel. I'm not sure if we want to stick with doing this in Zarf. It very well could make things slower, but would make us less prone to flakes