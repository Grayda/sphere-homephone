sphere-homephone
================

A Ninja Sphere driver that allows you, with the aid of a USB 56k modem, to get notified when your home phone is ringing. This driver also acts as a "call log" which you can visit later on and see who called, and when.

NOTE:
=====

Due to a [bug on the Ninja Sphere][1], this driver may not work. It may work on non-Sphere hardware, but with no LED matrix. Also, there are no (mobile) phone notifications, as the "notifications" channel on the Sphere is incomplete, in addition to the Sphere app not having such a feature (I think?)

How to use this driver
======================

Simply plug your USB 56k modem in (I use [this one][2]) and green reset your Sphere. When your modem is detected, a page will appear in the [Labs section][3] and you'll be able to see who called. If you don't have Caller ID turned on, the caller log will be empty

Compiling and Installing
========================

Copy `sphere-homephone` and `package.json` to `/data/sphere/user-autostart/drivers/sphere-homephone`. Create a folder called `images` in the same folder, then copy the two GIF files in there. They're not actually used in this version, but still good to copy over anyway.

To-Do
=====

- [ ] Get this to work with the Sphere
- [ ] Icons don't show on the Sphere (crash with nil pointer panic). Need to fix this
- [ ] Put in failsafes so if the port doesn't exist, don't crash, but show an error in the Labs
- [ ] Allow the user to change the COM port (if they're working on a BBB, Pi or Sphere, where the /dev entry might change)
- [ ] Add an "address book" so users can assign names to numbers (e.g. "Missed call from David" is better than "Missed call from +610412345678")

  [1]: https://discuss.ninjablocks.com/t/usb-device-s-not-recognised-on-the-sphere/2913/5?u=grayda
  [2]: http://goo.gl/gGZZne
  [3]: http://ninjasphere.local
