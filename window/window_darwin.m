#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

void requestAttention(int continuous) {
	@autoreleasepool {
		if (continuous)
			[NSApp requestUserAttention:NSCriticalRequest];
		else
			[NSApp requestUserAttention:NSInformationalRequest];
	}
}

