# Resourceful

Resourceful is an Xcode utility that helps avoid making mistakes when using `UIImage/NSImage imageNamed:`. That method takes a String, and loads an image from your app's bundle with that name. Unfortunately, there's no warning if you, say, typo the string.

Resourceful skirts this issue by inspecting your Xcode project for image assets, and creating an enum of all of their possible values. Thus,
```swift
let image = UIImage(named: "catsInSpace")
```

becomes
```swift
let image = Resource.Image.CatsInSpace.image
```

There are [many][0] [other][1] [libraries][2] that do this, but this one is mine. Why should you use it instead of anything else?
- It's written in Go, and really fast (important for a program that you'll be running as part of every Xcode build)
- It works for both iOS and OS X projects
- It can optionally generate warnings for you in places where you're still using `imageNamed` in your app.

## Installation

First, install Resourceful using homebrew:
```
brew install resourceful
```

Next, in your Xcode build phases settings, add a new Run Script Build Phase at the beginning of your build. For its contents, simply type
```
resourceful
```

In the "Input Files" section, put a link to your app's `Images.xcassets` directory. This will typically be:
```
$(SRCROOT)/Your_apps_name/Images.xcassets
```

In the "Output Files" section, specify where you'd like the code Resourceful generates to live. This will typically be:
```
$(SRCROOT)/Your_apps_name/Resourceful.swift
```

If you're wondering, you specify your files this way because Xcode can then intelligently avoid re-running this script unless the contents of your `Images.xcassets` changes.

Next, build your app. This will generate `Resourceful.swift` at the directory you specified.

Finally, add `Resourceful.swift` to your Xcode project by going to File -> Add Files to "Your app's name".

### Optional

As mentioned previously, Resourceful can also generate warnings in your Xcode project to encourage you to stop using imageNamed. To do this, add another Run Script Build Phase (it needs to be separate from the previous script, in case Xcode skips the previous one due to caching) with the contents:
```
resourceful warn
```

And you're done!

[0]: https://github.com/mac-cain13/R.swift
[1]: https://github.com/AliSoftware/SwiftGen
[2]: https://github.com/kaandedeoglu/Shark
