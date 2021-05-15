# File Format V2 WIP

New and improved file format.  
## Goals

 * Variable precision. If someone needs the highest precision possible, that's fine. If they don't, we can save them massive space.
 * Support Eye Tracking.
   * Another way of putting it: Allow people to build their own "thing" to record. Don't limit them to Unity GameObjects.
 * Support arbitrary binaries (to support things like audio), associate binaries with specific recording. Binaries need to also be able to have "lengths" associated with them, and when they start.
 * Support Groupings of Subjects
   * A great way of organizing things. Can group all "enemies" in one capture at a certain frame rate, and all things about the player in another grouping (Headset, right hand, and left hand)
 * Be able to support different "resolutions" in some form or fashion.
 * Support recording traditional animations.

## Definitions

* **Capture** - Data associated with a specific point in time
* **Capture Collection** - A collection of *Captures* around a specific type of data (think something like positional data). We can build many fundamental collections for people to use immediately: Vector3, float, int, dictionary, etc.
* **Recording** - A collection of capture collections and recordings.

## Binary File Spec

### Common Binary Patterns

Common patterns you'll see throughout the document.

* [varint] - variable sized uint64. Only takes up 1 byte if only one byte is needed.
* [string] - varint representing number of characters in string, followed by string itself.

